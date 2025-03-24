/* eslint-disable @typescript-eslint/no-unused-vars */
import { execSync } from 'node:child_process';
import { basename, dirname, join } from 'node:path';
import { NxJsonConfiguration, ProjectGraph, Tree } from '@nx/devkit';
import { ManifestData, VersionActions } from 'nx/release';

// NOTE: LIMITATION: This assumes the name of the package from the last part of the path.
// So having two package with the same name in different directories will cause issues.

/**
 * Implementation of version actions for Go projects.
 * This class handles versioning operations for Go modules using Git tags.
 */
export default class GoVersionActions extends VersionActions {
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

    // Get the full project root directory
    const projectRoot = dirname(sourcePath);

    // Get a more readable name from the directory
    const name = basename(projectRoot);

    const dependencies = {};

    // Get current version from Git tags based on Nx release configuration
    const currentVersion = await this.getCurrentVersionFromGit(tree);

    return {
      name,
      currentVersion,
      dependencies,
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
    // Get current version from Git tags based on Nx release configuration
    const currentVersion = await this.getCurrentVersionFromGit(tree);

    return currentVersion;
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
      const sourcePath = this.getSourceManifestPath();
      const content = tree.read(sourcePath, 'utf-8');
      const moduleNameMatch = content.match(/module\s+([^\s]+)/);
      const moduleName = moduleNameMatch ? moduleNameMatch[1] : '';

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

  async getCurrentVersionOfDependency(
    tree: Tree,
    projectGraph: ProjectGraph,
    dependencyProjectName: string
  ): Promise<{
    currentVersion: string | null;
    dependencyCollection: string | null;
  }> {
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
    // eslint-disable-next-line @typescript-eslint/no-empty-function
  ): Promise<void> {}

  /**
   * Updates go.mod content with new dependency versions
   */
}
