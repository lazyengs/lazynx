import {
  CreateNodesContextV2,
  CreateNodesV2,
  ProjectConfiguration,
  TargetConfiguration,
  createNodesFromFiles,
  joinPathFragments,
  workspaceRoot,
} from '@nx/devkit';
import { basename, dirname, join } from 'path';
import { existsSync, readFileSync } from 'fs';

// File glob to find all Go projects
const goModGlob = '**/go.mod';

// Expected format of the plugin options defined in nx.json
export interface GoPluginOptions {
  buildTargetName?: string;
  testTargetName?: string;
  runTargetName?: string;
  tidyTargetName?: string;
  lintTargetName?: string;
  tagName?: string;
}

// Entry function that Nx calls to modify the graph
export const createNodesV2: CreateNodesV2<GoPluginOptions> = [
  goModGlob,
  (configFiles, options, context) => {
    return createNodesFromFiles(
      (configFile, options, context) =>
        createNodesInternal(configFile, options, context),
      configFiles,
      options,
      context
    );
  },
];

function createNodesInternal(
  configFilePath: string,
  options: GoPluginOptions,
  context: CreateNodesContextV2
) {
  // Get the full project root directory
  const projectRoot = dirname(configFilePath);

  // Get a more readable name from the directory
  const projectName = basename(projectRoot);

  // For better UX, set default target names if not provided
  const buildTargetName = options.buildTargetName || 'build';
  const testTargetName = options.testTargetName || 'test';
  const runTargetName = options.runTargetName || 'serve';
  const tidyTargetName = options.tidyTargetName || 'tidy';
  const lintTargetName = options.lintTargetName || 'lint';

  // Detect if this is an application or a library
  const isApplication = hasMainPackage(projectRoot);
  const projectType = isApplication ? 'application' : 'library';

  // Initialize targets object
  const targets: Record<string, TargetConfiguration> = {};

  // Common test target - available for both apps and libraries
  targets[testTargetName] = {
    executor: 'nx:run-commands',
    options: {
      command: 'go test ./...',
      cwd: projectRoot,
    },
    cache: true,
    inputs: [
      '{projectRoot}/go.mod',
      '{projectRoot}/go.sum',
      joinPathFragments('{projectRoot}', '**', '*.go'),
      {
        externalDependencies: ['go'],
      },
    ],
    outputs: [],
  };

  // Tidy target - available for both apps and libraries
  targets[tidyTargetName] = {
    executor: 'nx:run-commands',
    options: {
      command: 'go mod tidy',
      cwd: projectRoot,
    },
  };

  // Lint target - available for both apps and libraries
  targets[lintTargetName] = {
    executor: 'nx:run-commands',
    options: {
      command: 'golangci-lint run',
      cwd: projectRoot,
    },
  };

  // Build and run targets - only for applications
  if (isApplication) {
    targets[buildTargetName] = {
      executor: 'nx:run-commands',
      options: {
        command: 'go build .',
        cwd: projectRoot,
      },
      cache: true,
      inputs: [
        '{projectRoot}/go.mod',
        '{projectRoot}/go.sum',
        joinPathFragments('{projectRoot}', '**', '*.go'),
        {
          externalDependencies: ['go'],
        },
      ],
      outputs: ['{projectRoot}/' + projectName],
    };

    targets[runTargetName] = {
      executor: 'nx:run-commands',
      options: {
        command: 'go run .',
        cwd: projectRoot,
      },
    };
  }

  // Create the project configuration
  const projectConfig: ProjectConfiguration & { root: string } = {
    name: projectName,
    root: projectRoot,
    sourceRoot: projectRoot,
    projectType,
    targets,
    tags: options.tagName ? [options.tagName] : [],
  };

  return {
    projects: {
      [projectRoot]: projectConfig,
    },
  };
}

/**
 * Determines if a Go module contains a main package,
 * indicating it is an application rather than a library.
 *
 * @param projectRoot The root directory of the project
 * @returns True if the project is an application, false otherwise
 */
function hasMainPackage(projectRoot: string): boolean {
  try {
    // Check if main.go exists in the project root
    const mainGoPath = join(workspaceRoot, projectRoot, 'main.go');
    if (existsSync(mainGoPath)) {
      const content = readFileSync(mainGoPath, 'utf-8');
      return content.includes('package main') && content.includes('func main(');
    }

    // Check for cmd directory structure (common Go pattern)
    const cmdDirPath = join(workspaceRoot, projectRoot, 'cmd');
    if (existsSync(cmdDirPath)) {
      return true;
    }

    return false;
  } catch (error) {
    // If there's any error, default to treating it as a library
    return false;
  }
}
