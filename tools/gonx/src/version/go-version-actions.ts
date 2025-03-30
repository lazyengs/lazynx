/* eslint-disable @typescript-eslint/no-unused-vars */
import { ProjectGraph, Tree } from '@nx/devkit';
import * as path from 'node:path';
import { VersionActions } from 'nx/release';
import { NxReleaseVersionV2Configuration } from 'nx/src/config/nx-json';

type ProxyGolangOrgVersionResponse = {
  Version: string;
  Time: string;
  Origin: {
    VCS: 'git';
    URL: string;
    Ref: string;
    Hash: string;
  };
};

const MANIFEST_FILENAME = 'go.mod';

// NOTE: LIMITATION: This assumes the name of the package from the last part of the path.
// So having two package with the same name in different directories will cause issues.

/**
 * Implementation of version actions for Go projects.
 * This class handles versioning operations for Go modules using Git tags.
 */
export default class GoVersionActions extends VersionActions {
  validManifestFilenames: string[] = [MANIFEST_FILENAME];

  // go.mod don't contain the version of the package, so we need to get the version from Git tags or from the registry
  readCurrentVersionFromSourceManifest(
    tree: Tree
  ): Promise<{ currentVersion: string; manifestPath: string } | null> {
    return null;
  }

  async readCurrentVersionFromRegistry(
    tree: Tree,
    currentVersionResolverMetadata: NxReleaseVersionV2Configuration['currentVersionResolverMetadata']
  ): Promise<{ currentVersion: string | null; logText: string } | null> {
    try {
      const manifestPath = path.join(
        this.projectGraphNode.data.sourceRoot,
        MANIFEST_FILENAME
      );
      const content = tree.read(manifestPath, 'utf-8');
      const moduleNameMatch = content.match(/module\s+([^\s]+)/);
      const moduleName = moduleNameMatch ? moduleNameMatch[1] : '';

      const result = await fetch(
        `https://proxy.golang.org/${encodeURIComponent(moduleName)}/@latest`
      );

      if (result?.ok) {
        const response = (await result.json()) as ProxyGolangOrgVersionResponse;
        // Get the latest version from the list
        const latestVersion = response.Version;

        return {
          currentVersion: latestVersion,
          logText: `Retrieved version ${latestVersion} from proxy.golang.org for ${moduleName}`,
        };
      }
    } catch (error) {
      console.error('Error checking proxy.golang.org:', error);

      return null;
    }
  }

  readCurrentVersionOfDependency(
    tree: Tree,
    projectGraph: ProjectGraph,
    dependencyProjectName: string
  ): Promise<{
    currentVersion: string | null;
    dependencyCollection: string | null;
  }> {
    throw new Error('Method not implemented.');
  }

  isLocalDependencyProtocol(versionSpecifier: string): Promise<boolean> {
    throw new Error('Method not implemented.');
  }

  async updateProjectVersion(
    tree: Tree,
    newVersion: string
  ): Promise<string[]> {
    // We do nothing on go projects by default
    return [];
  }

  async updateProjectDependencies(
    tree: Tree,
    projectGraph: ProjectGraph,
    dependenciesToUpdate: Record<string, string>
  ): Promise<string[]> {
    return [];
  }
}
