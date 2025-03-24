import { execSync } from 'child_process';
import { dirname, join } from 'path';
import { NxJsonConfiguration, ProjectGraph, Tree } from '@nx/devkit';
import { ManifestData, VersionActions } from 'nx/release';

// NOTE: LIMITATION: This assumes the name of the package from the last part of the path.
// So having two package with the same name in different directories will cause issues.

/**
 * Implementation of version actions for Go projects.
 * This class handles versioning operations for Go modules using Git tags.
 */
export class GoVersionActions extends VersionActions {
  /**
   * The filename for Go module manifest
   */
  manifestFilename = 'go.mod';

  /**
   * Reads the Go module manifest data from the source
   * @param tree The file tree
   * @returns ManifestData with Go module information
   */
  async readSourceManifestData(tree: Tree): Promise<ManifestData> {
    const sourcePath = this.getSourceManifestPath();
    if (!tree.exists(sourcePath)) {
      throw new Error(
        `Unable to read the go.mod file at ${sourcePath}, please ensure that the file exists and is valid`
      );
    }

    const content = tree.read(sourcePath, 'utf-8');
    const moduleNameMatch = content.match(/module\s+([^\s]+)/);
    const moduleName = moduleNameMatch ? moduleNameMatch[1] : '';
    const name = moduleName.split('/').pop() || moduleName;
    const dependencies = this.parseDependencies(sourcePath, tree);

    // Get current version from Git tags based on Nx release configuration
    const currentVersion = await this.getCurrentVersionFromGit(tree);

    return {
      name,
      currentVersion,
      dependencies,
    };
  }

  /**
   * Parses dependencies from a go.mod file and source code imports
   * @param sourcePath The path to the go.mod file
   * @returns Record of dependencies
   */
  private parseDependencies(sourcePath: string, tree: Tree) {
    // Extract dependencies
    const dependencies: Record<
      string,
      Record<string, { resolvedVersion: string; rawVersionSpec: string }>
    > = {
      dependencies: {},
      devDependencies: {},
    };

    try {
      // Read the go.mod file
      const goModContent = tree.read(sourcePath, 'utf-8');
      if (!goModContent) {
        return dependencies;
      }

      // Parse require blocks
      const requireSections = goModContent.match(
        /require\s+\(\s*([^)]+)\s*\)/gs
      );
      if (requireSections) {
        for (const section of requireSections) {
          const deps = section.match(/\s*([^\s]+)\s+([^\s]+)/g);
          if (deps) {
            for (const dep of deps) {
              const [name, version] = dep.trim().split(/\s+/);
              if (name && version) {
                dependencies.dependencies[name] = {
                  resolvedVersion: version,
                  rawVersionSpec: version,
                };
              }
            }
          }
        }
      }

      // Parse single-line requires
      const singleRequires = goModContent.match(
        /require\s+([^\s]+)\s+([^\s]+)/g
      );
      if (singleRequires) {
        for (const req of singleRequires) {
          const parts = req.split(/\s+/);
          if (parts.length >= 3) {
            const name = parts[1];
            const version = parts[2];
            dependencies.dependencies[name] = {
              resolvedVersion: version,
              rawVersionSpec: version,
            };
          }
        }
      }

      // Parse replace directives
      const replaceSections = goModContent.match(
        /replace\s+([^\s]+)\s+=>\s+([^\s]+)(\s+([^\s]+))?/g
      );
      if (replaceSections) {
        for (const replaceSection of replaceSections) {
          const match = replaceSection.match(
            /replace\s+([^\s]+)\s+=>\s+([^\s]+)(\s+([^\s]+))?/
          );
          if (match) {
            const [, moduleName, replacePath, , version] = match;

            // If the replacement is a local path (starts with ./ or ../ or absolute), consider it a local dependency
            if (
              replacePath.startsWith('./') ||
              replacePath.startsWith('../') ||
              replacePath.startsWith('/')
            ) {
              // Local replacement - could be in workspace
              if (version) {
                dependencies.dependencies[moduleName] = {
                  resolvedVersion: version,
                  rawVersionSpec: `${replacePath} ${version}`,
                };
              } else {
                dependencies.dependencies[moduleName] = {
                  resolvedVersion: 'local',
                  rawVersionSpec: replacePath,
                };
              }
            } else {
              // Remote replacement
              if (version) {
                dependencies.dependencies[moduleName] = {
                  resolvedVersion: version,
                  rawVersionSpec: `${replacePath} ${version}`,
                };
              }
            }
          }
        }
      }

      // Look for imports in Go source files
      this.addDependenciesFromGoImports(dependencies, tree);
    } catch (error) {
      console.error(`Error parsing dependencies from ${sourcePath}:`, error);
    }

    return dependencies;
  }

  /**
   * Find imported packages in Go source files and add them to dependencies
   */
  private addDependenciesFromGoImports(
    dependencies: Record<
      string,
      Record<string, { resolvedVersion: string; rawVersionSpec: string }>
    >,
    tree: Tree
  ): void {
    try {
      const projectRoot = this.projectGraphNode.data.root;

      // Find all Go files in this project
      const findGoFilesCmd = `find ${projectRoot} -name "*.go" -type f`;
      const goFiles = execSync(findGoFilesCmd, {
        encoding: 'utf-8',
        cwd: process.cwd(),
      })
        .trim()
        .split('\n')
        .filter(Boolean);

      // Regular expression for Go imports
      const importRegex = /import\s+(?:\w+\s+)?"([^"]+)"|import\s+\(([^)]*)\)/g;
      const singleImportRegex = /"([^"]+)"/g;

      for (const filePath of goFiles) {
        if (!tree.exists(filePath)) continue;

        const fileContent = tree.read(filePath, 'utf-8');

        // Extract imports from each file
        const importMatches = Array.from(fileContent.matchAll(importRegex));

        for (const match of importMatches) {
          if (match[1]) {
            // Single import
            this.addImportToDependencies(match[1], dependencies);
          } else if (match[2]) {
            // Import block
            const innerMatches = Array.from(
              match[2].matchAll(singleImportRegex)
            );
            for (const innerMatch of innerMatches) {
              this.addImportToDependencies(innerMatch[1], dependencies);
            }
          }
        }
      }
    } catch (error) {
      console.error('Error parsing Go imports:', error);
    }
  }

  /**
   * Process a single import to see if it should be added to dependencies
   */
  private addImportToDependencies(
    importPath: string,
    dependencies: Record<
      string,
      Record<string, { resolvedVersion: string; rawVersionSpec: string }>
    >
  ): void {
    // Skip standard library imports
    if (!importPath.includes('.') && !importPath.includes('/')) {
      return;
    }

    // Skip relative imports
    if (importPath.startsWith('./') || importPath.startsWith('../')) {
      return;
    }

    // Skip already tracked dependencies
    if (dependencies.dependencies[importPath]) {
      return;
    }

    // Find potential module name from import path
    const potentialModule = importPath;

    // Check if this import might be from a workspace module but not in go.mod
    // This could happen in Go workspaces where local modules don't need requires
    // For now, we'll add it as an implicit dependency with a default version
    dependencies.dependencies[potentialModule] = {
      resolvedVersion: 'v0.0.0',
      rawVersionSpec: 'v0.0.0',
    };
  }

  /**
   * Gets the release configuration from nx.json to determine fixed/independent mode and tag patterns
   * @param tree The file tree
   */
  private getNxReleaseConfig(tree: Tree): {
    projectsRelationship: 'fixed' | 'independent';
    releaseTagPattern: string;
    projectReleaseTagPattern: string | null;
  } {
    const nxJsonPath = 'nx.json';
    if (!tree.exists(nxJsonPath)) {
      throw new Error('nx.json not found');
    }

    const nxJson = JSON.parse(
      tree.read(nxJsonPath, 'utf-8')
    ) as NxJsonConfiguration;
    const releaseConfig = nxJson.release || {};

    // Determine relationship mode (fixed or independent)
    const projectsRelationship = releaseConfig.projectsRelationship ?? 'fixed';

    // Get the default release tag pattern
    const defaultReleaseTagPattern =
      releaseConfig.releaseTagPattern ??
      (projectsRelationship === 'fixed'
        ? 'v{version}'
        : '{projectName}@{version}');

    // Check if the project is in a specific group with a custom tag pattern
    let projectReleaseTagPattern = null;
    const projectName = this.projectGraphNode.name;

    if (releaseConfig.groups) {
      // Find the group containing this project
      for (const [, groupConfig] of Object.entries(releaseConfig.groups)) {
        const groupProjects =
          typeof groupConfig.projects === 'string'
            ? [groupConfig.projects]
            : groupConfig.projects;

        if (groupProjects.includes(projectName)) {
          // Use the group-specific tag pattern if available
          projectReleaseTagPattern = groupConfig.releaseTagPattern || null;
          break;
        }
      }
    }

    return {
      projectsRelationship,
      releaseTagPattern: defaultReleaseTagPattern,
      projectReleaseTagPattern,
    };
  }

  /**
   * Get the current version from Git tags based on Nx release configuration
   * @returns The current version from Git tags or v0.0.0 if none
   */
  private async getCurrentVersionFromGit(tree: Tree): Promise<string> {
    const {
      projectsRelationship,
      releaseTagPattern,
      projectReleaseTagPattern,
    } = this.getNxReleaseConfig(tree);

    // Determine which tag pattern to use
    const tagPatternToUse = projectReleaseTagPattern || releaseTagPattern;

    // Create the Git match pattern based on the release configuration
    let gitMatchPattern: string;

    if (projectsRelationship === 'independent' || projectReleaseTagPattern) {
      // For independent projects or projects with custom patterns, we need to replace {projectName} with the actual name
      gitMatchPattern = tagPatternToUse
        .replace('{projectName}', this.projectGraphNode.name)
        .replace('{version}', '*');
    } else {
      // For fixed projects with the default pattern, replace {version} with a wildcard
      gitMatchPattern = tagPatternToUse.replace('{version}', '*');
    }

    // Get latest tag for this module using the appropriate pattern
    const gitCommand = `git describe --tags --abbrev=0 --match "${gitMatchPattern}" 2>/dev/null || echo "v0.0.0"`;

    const result = execSync(gitCommand, {
      encoding: 'utf-8',
      cwd: process.cwd(),
    }).trim();

    if (result === '' || result === 'v0.0.0') {
      return 'v0.0.0';
    }

    // Extract the version from the tag based on the pattern
    return this.extractVersionFromTag(result, tagPatternToUse);
  }

  /**
   * Extracts the version component from a Git tag based on the release tag pattern
   * @param tag The full Git tag
   * @param pattern The pattern used to generate the tag
   * @returns The extracted version
   */
  private extractVersionFromTag(tag: string, pattern: string): string {
    try {
      if (pattern === 'v{version}') {
        // Remove v prefix if present
        return tag.startsWith('v') ? tag.slice(1) : tag;
      }

      if (pattern === '{projectName}@{version}') {
        // Extract version after @
        const parts = tag.split('@');
        if (parts.length >= 2) {
          return parts[1];
        }
      }

      if (pattern === '{projectName}@v{version}') {
        // Extract version after @ and remove v prefix if present
        const parts = tag.split('@');
        if (parts.length >= 2) {
          return parts[1].startsWith('v') ? parts[1].slice(1) : parts[1];
        }
      }

      // For other custom patterns, try to extract based on semver detection
      const semverMatch = tag.match(
        /(\d+\.\d+\.\d+(?:-[0-9A-Za-z-.]+)?(?:\+[0-9A-Za-z-.]+)?)/
      );
      if (semverMatch) {
        return semverMatch[1];
      }

      // If we can't extract, return the original tag
      return tag;
    } catch {
      // If any error occurs, return the original tag
      return tag;
    }
  }

  /**
   * Read the current version from Git tags
   * @param tree The file tree
   * @returns The current version
   */
  async readCurrentVersionFromSourceManifest(tree: Tree): Promise<string> {
    const sourceData = await this.readCachedSourceManifestData(tree);
    return sourceData.currentVersion;
  }

  /**
   * Read the current version from a remote registry
   * This attempts to check proxy.golang.org or pkg.go.dev for version info
   * @param tree The file tree
   * @param currentVersionResolverMetadata Optional metadata for resolving versions
   * @returns The current version from registry and log text
   */
  async readCurrentVersionFromRegistry(
    tree: Tree
  ): Promise<{ currentVersion: string; logText: string }> {
    try {
      const sourceData = await this.readCachedSourceManifestData(tree);
      const moduleName = sourceData.name;

      // Try to get the latest published version from proxy.golang.org
      const result = execSync(
        `curl -s https://proxy.golang.org/${encodeURIComponent(
          moduleName
        )}/@v/list`,
        {
          encoding: 'utf-8',
        }
      ).trim();

      if (result) {
        // Get the latest version from the list
        const versions = result.split('\n').filter((v) => v.startsWith('v'));
        if (versions.length > 0) {
          // Sort versions semantically and get the latest
          versions.sort((a, b) => {
            return a.localeCompare(b, undefined, {
              numeric: true,
              sensitivity: 'base',
            });
          });
          const latestVersion = versions[versions.length - 1];

          return {
            currentVersion: latestVersion,
            logText: `Retrieved version ${latestVersion} from proxy.golang.org for ${moduleName}`,
          };
        }
      }

      // If proxy doesn't have info, fall back to Git tags
      const currentVersion = await this.getCurrentVersionFromGit(tree);
      return {
        currentVersion,
        logText: `Using version ${currentVersion} from Git tags (proxy.golang.org had no information)`,
      };
    } catch (error) {
      console.error('Error checking proxy.golang.org:', error);
      // If we can't get the version from proxy, fall back to Git tags
      const currentVersion = await this.getCurrentVersionFromGit(tree);
      return {
        currentVersion,
        logText: `Using version ${currentVersion} from Git tags (proxy.golang.org check failed)`,
      };
    }
  }

  /**
   * Get the current version of a dependency from the go.mod file
   * @param tree The file tree
   * @param projectGraph The project graph
   * @param dependencyProjectName The dependency project name
   * @returns The current version and dependency collection
   */
  async getCurrentVersionOfDependency(
    tree: Tree,
    projectGraph: ProjectGraph,
    dependencyProjectName: string
  ): Promise<{
    currentVersion: string | null;
    dependencyCollection: string | null;
  }> {
    const dependencyNode = projectGraph.nodes[dependencyProjectName];

    if (!dependencyNode) {
      return { currentVersion: null, dependencyCollection: null };
    }

    // Get the module name of the dependency
    const dependencyModulePath = join(dependencyNode.data.root, 'go.mod');
    if (!tree.exists(dependencyModulePath)) {
      return { currentVersion: null, dependencyCollection: null };
    }

    const dependencyModContent = tree.read(dependencyModulePath, 'utf-8');
    const moduleMatch = dependencyModContent.match(/module\s+([^\s]+)/);
    if (!moduleMatch) {
      return { currentVersion: null, dependencyCollection: null };
    }

    const moduleName = moduleMatch[1];

    // Parse this project's go.mod and source files to find dependencies
    // We'll do this directly instead of using the cached manifest data to avoid circular dependencies
    const projectRoot = this.projectGraphNode.data.root;
    const goModPath = join(projectRoot, 'go.mod');

    // First check for explicit dependencies in go.mod
    if (tree.exists(goModPath)) {
      const goModContent = tree.read(goModPath, 'utf-8');

      // Check for dependency in require block
      const requireBlockMatch = goModContent.match(
        new RegExp(
          `require\\s+\\([^)]*${escapeRegExp(moduleName)}\\s+([^\\s]+)`
        )
      );

      if (requireBlockMatch && requireBlockMatch[1]) {
        return {
          currentVersion: requireBlockMatch[1],
          dependencyCollection: 'dependencies',
        };
      }

      // Check for dependency in single-line require
      const singleRequireMatch = goModContent.match(
        new RegExp(`require\\s+${escapeRegExp(moduleName)}\\s+([^\\s]+)`)
      );

      if (singleRequireMatch && singleRequireMatch[1]) {
        return {
          currentVersion: singleRequireMatch[1],
          dependencyCollection: 'dependencies',
        };
      }
    }

    // If not found in go.mod explicitly, check for imports in Go files
    const dependencies = this.findWorkspaceDependencies(
      tree,
      projectGraph,
      projectRoot
    );

    if (dependencies.has(dependencyProjectName)) {
      return {
        currentVersion:
          dependencies.get(dependencyProjectName)?.version || null,
        // We'll use 'dependencies' as the collection for imports found in Go files
        dependencyCollection: 'dependencies',
      };
    }

    return { currentVersion: null, dependencyCollection: null };
  }

  /**
   * Check if a version specifier uses a local dependency protocol
   * @param versionSpecifier The version specifier to check
   * @returns True if it's a local dependency
   */
  isLocalDependencyProtocol(versionSpecifier: string): boolean {
    // In Go, local replacements would be in replace directives
    // A simplified check for local paths
    return (
      versionSpecifier.startsWith('./') ||
      versionSpecifier.startsWith('../') ||
      versionSpecifier.startsWith('/') ||
      /v0\.0\.0-\d{14}-[a-f0-9]{12}/.test(versionSpecifier) // Pseudo-versions often used for local deps
    );
  }

  /**
   * Parse Go files to find dependencies between modules directly from file content
   * Similar to the logic in create-dependencies.ts
   * @param tree The file tree
   * @param projectGraph The project graph
   * @param thisProjectRoot The root directory of the current project
   * @returns Map of module dependencies
   */
  private findWorkspaceDependencies(
    tree: Tree,
    projectGraph: ProjectGraph,
    thisProjectRoot: string
  ): Map<string, { moduleName: string; version: string | null }> {
    const dependencies = new Map<
      string,
      { moduleName: string; version: string | null }
    >();

    // First read this module's go.mod to get explicit dependencies
    const goModPath = join(thisProjectRoot, 'go.mod');
    const goModDependencies: Record<string, string> = {};

    if (tree.exists(goModPath)) {
      const goModContent = tree.read(goModPath, 'utf-8');

      // Extract dependencies from require sections
      const requireSections = goModContent.match(
        /require\s+\(\s*([^)]+)\s*\)/g
      );
      if (requireSections) {
        for (const section of requireSections) {
          const deps = section.match(/\s*([^\s]+)\s+([^\s]+)/g);
          if (deps) {
            for (const dep of deps) {
              const [name, version] = dep.trim().split(/\s+/);
              if (name && version) {
                goModDependencies[name] = version;
              }
            }
          }
        }
      }

      // Also check for single-line require statements
      const singleRequires = goModContent.match(
        /require\s+([^\s]+)\s+([^\s]+)/g
      );
      if (singleRequires) {
        for (const req of singleRequires) {
          const parts = req.split(/\s+/);
          if (parts.length >= 3) {
            const name = parts[1];
            const version = parts[2];
            goModDependencies[name] = version;
          }
        }
      }
    }

    // Build a map of all module names to project names in the workspace
    const moduleToProjectMap = new Map<string, string>();

    // Find all Go modules in the workspace
    for (const [projectName, node] of Object.entries(projectGraph.nodes)) {
      const modPath = join(node.data.root, 'go.mod');
      if (tree.exists(modPath)) {
        const modContent = tree.read(modPath, 'utf-8');
        const moduleMatch = modContent.match(/module\s+([^\s]+)/);
        if (moduleMatch && moduleMatch[1]) {
          moduleToProjectMap.set(moduleMatch[1], projectName);
        }
      }
    }

    // Regular expression for Go imports
    const importRegex = /import\s+(?:\w+\s+)?"([^"]+)"|import\s+\(([^)]*)\)/g;
    const singleImportRegex = /"([^"]+)"/g;

    // Find all Go files in this project
    try {
      const findGoFilesCmd = `find ${thisProjectRoot} -name "*.go" -type f`;
      const goFiles = execSync(findGoFilesCmd, {
        encoding: 'utf-8',
        cwd: process.cwd(),
      })
        .trim()
        .split('\n')
        .filter(Boolean);

      // Process each Go file
      for (const filePath of goFiles) {
        if (!tree.exists(filePath)) continue;

        const fileContent = tree.read(filePath, 'utf-8');

        // Extract imports from each file
        const importMatches = Array.from(fileContent.matchAll(importRegex));

        for (const match of importMatches) {
          if (match[1]) {
            // Single import
            this.checkWorkspaceImport(
              match[1],
              moduleToProjectMap,
              dependencies,
              goModDependencies
            );
          } else if (match[2]) {
            // Import block
            const innerMatches = Array.from(
              match[2].matchAll(singleImportRegex)
            );
            for (const innerMatch of innerMatches) {
              this.checkWorkspaceImport(
                innerMatch[1],
                moduleToProjectMap,
                dependencies,
                goModDependencies
              );
            }
          }
        }
      }
    } catch (error) {
      // If there's an error finding Go files, log it and continue
      console.error('Error parsing Go files for dependencies:', error);
    }

    return dependencies;
  }

  /**
   * Check if an import matches a workspace module and add it to dependencies
   */
  private checkWorkspaceImport(
    importPath: string,
    moduleToProjectMap: Map<string, string>,
    dependencies: Map<string, { moduleName: string; version: string | null }>,
    goModDependencies: Record<string, string>
  ): void {
    // Check if the import matches any module in our workspace
    for (const [moduleName, projectName] of moduleToProjectMap.entries()) {
      if (
        importPath === moduleName ||
        importPath.startsWith(moduleName + '/')
      ) {
        // Found a match - this import is from one of our workspace modules

        // Check if we have an explicit version in go.mod
        const version = goModDependencies[moduleName] || null;

        // Add to dependencies map if not already present
        if (!dependencies.has(projectName)) {
          dependencies.set(projectName, { moduleName, version });
        }
        break;
      }
    }
  }

  /**
   * Write a new version by creating a Git tag
   * This doesn't modify the go.mod file directly, but creates a new Git tag
   * @param tree The file tree
   * @param newVersion The new version to tag
   */
  async writeVersionToManifests(tree: Tree, newVersion: string): Promise<void> {
    // For Go modules, we don't modify the go.mod file to update the version
    // Instead, we need to create a Git tag which will be reflected when the module is published

    // Since we can't create Git tags directly within the Tree API,
    // we need to add a build script to be executed after files are written
    // This will be handled by the afterAllProjectsVersionedCallback
    // which will create the Git tag when it runs

    // Store the new version in a temporary file that the callback can read
    const versionFilePath = join(
      this.projectGraphNode.data.root,
      '.version-temp'
    );
    tree.write(versionFilePath, newVersion);

    // For now, we'll update a VERSION file if it exists (common practice in Go projects)
    const versionFileExists = tree.exists(
      join(this.projectGraphNode.data.root, 'VERSION')
    );
    if (versionFileExists) {
      tree.write(join(this.projectGraphNode.data.root, 'VERSION'), newVersion);
    }
  }

  async updateDependencies(
    tree: Tree,
    projectGraph: ProjectGraph,
    dependenciesToUpdate: Record<string, string>
  ): Promise<void> {
    // For each dependency to update, we need to:
    // 1. Find its module name
    // 2. Update the version in each go.mod file using go mod edit

    const dependencyModuleNames: Record<string, string> = {};
    const dependencyVersionUpdates: Record<string, string> = {};

    // First, get all module names for the dependencies
    for (const dependencyName of Object.keys(dependenciesToUpdate)) {
      const dependencyNode = projectGraph.nodes[dependencyName];
      if (!dependencyNode) continue;

      const dependencyModPath = join(dependencyNode.data.root, 'go.mod');
      if (!tree.exists(dependencyModPath)) continue;

      const modContent = tree.read(dependencyModPath, 'utf-8');
      const moduleMatch = modContent.match(/module\s+([^\s]+)/);
      if (moduleMatch) {
        const moduleName = moduleMatch[1];
        dependencyModuleNames[dependencyName] = moduleName;
        dependencyVersionUpdates[moduleName] =
          dependenciesToUpdate[dependencyName];
      }
    }

    // Now update each manifest file
    for (const manifestPath of this.manifestsToUpdate) {
      if (!tree.exists(manifestPath)) continue;

      const content = tree.read(manifestPath, 'utf-8');
      const moduleDir = dirname(manifestPath);

      // Parse dependencies from imports to check if we need to add new requires
      const directDependencies = this.findWorkspaceDependencies(
        tree,
        projectGraph,
        moduleDir
      );
      const updatedContent = this.updateGoModContent(
        content,
        dependencyVersionUpdates,
        directDependencies,
        dependencyModuleNames
      );

      tree.write(manifestPath, updatedContent);
    }

    // Also create a temporary file with version updates for the callback
    const updateFilePath = join(
      this.projectGraphNode.data.root,
      '.dependency-updates-temp'
    );
    tree.write(updateFilePath, JSON.stringify(dependencyVersionUpdates));
  }

  /**
   * Updates go.mod content with new dependency versions
   */
  private updateGoModContent(
    content: string,
    dependencyVersionUpdates: Record<string, string>,
    directDependencies: Map<
      string,
      { moduleName: string; version: string | null }
    >,
    dependencyNameMap: Record<string, string>
  ): string {
    let updatedContent = content;
    const hasRequireBlock = /require\s+\(/.test(content);
    const requireBlockEndMatch = content.match(/require\s+\([^)]*\)/s);

    // Track which dependencies we've updated
    const updatedModules = new Set<string>();

    // First update existing dependencies in requires
    for (const [moduleName, newVersion] of Object.entries(
      dependencyVersionUpdates
    )) {
      // Update in require blocks
      const requireBlockPattern = new RegExp(
        `(require\\s+\\([^)]*${escapeRegExp(
          moduleName
        )}\\s+)([^\\s]+)([^)]*\\))`,
        'g'
      );

      // Check if we updated anything
      const contentBeforeUpdate = updatedContent;
      updatedContent = updatedContent.replace(
        requireBlockPattern,
        `$1${newVersion}$3`
      );

      if (contentBeforeUpdate !== updatedContent) {
        updatedModules.add(moduleName);
      } else {
        // Check for single-line requires
        const singleRequirePattern = new RegExp(
          `(require\\s+${escapeRegExp(moduleName)}\\s+)([^\\s]+)`,
          'g'
        );
        const contentBeforeSingleUpdate = updatedContent;
        updatedContent = updatedContent.replace(
          singleRequirePattern,
          `$1${newVersion}`
        );

        if (contentBeforeSingleUpdate !== updatedContent) {
          updatedModules.add(moduleName);
        }
      }
    }

    // Now check for dependencies we found in imports but aren't in go.mod yet
    // We'll need to add requires for these
    const reverseMap = Object.entries(dependencyNameMap).reduce(
      (acc, [projectName, moduleName]) => {
        acc[moduleName] = projectName;
        return acc;
      },
      {} as Record<string, string>
    );

    for (const [, { moduleName }] of directDependencies.entries()) {
      // Skip if we don't have an update for this module
      if (!reverseMap[moduleName] || !dependencyVersionUpdates[moduleName])
        continue;

      // Skip if already updated
      if (updatedModules.has(moduleName)) continue;

      // Module needs to be added to requires
      const newVersion = dependencyVersionUpdates[moduleName];

      if (hasRequireBlock && requireBlockEndMatch) {
        // Add to existing require block
        const requireBlockEnd = requireBlockEndMatch[0];
        const insertPoint = requireBlockEnd.lastIndexOf(')');
        const newRequireContent =
          requireBlockEnd.slice(0, insertPoint) +
          `\t${moduleName} ${newVersion}\n` +
          requireBlockEnd.slice(insertPoint);

        updatedContent = updatedContent.replace(
          requireBlockEnd,
          newRequireContent
        );
      } else {
        // Add a new require line
        if (hasRequireBlock) {
          // There's a require block but we couldn't find its end, just append
          updatedContent += `\nrequire ${moduleName} ${newVersion}\n`;
        } else {
          // No require block found, add a new one after the module line
          const moduleLineMatch = updatedContent.match(/module\s+[^\s]+/);
          if (moduleLineMatch) {
            const moduleLine = moduleLineMatch[0];
            updatedContent = updatedContent.replace(
              moduleLine,
              `${moduleLine}\n\nrequire ${moduleName} ${newVersion}`
            );
          } else {
            // No module line found, just append
            updatedContent += `\nrequire ${moduleName} ${newVersion}\n`;
          }
        }
      }

      updatedModules.add(moduleName);
    }

    return updatedContent;
  }
}

/**
 * Helper function to escape special regex characters in strings
 */
function escapeRegExp(string: string): string {
  return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

/**
 * Export the Go version actions instance
 */
export const goVersionActions = GoVersionActions;
