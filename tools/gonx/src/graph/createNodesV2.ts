import {
  CreateNodesContextV2,
  CreateNodesV2,
  ProjectConfiguration,
  TargetConfiguration,
  createNodesFromFiles,
  joinPathFragments,
} from '@nx/devkit';
import { basename, dirname } from 'path';

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

  // Build target configuration
  const buildTarget: TargetConfiguration = {
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

  // Test target configuration
  const testTarget: TargetConfiguration = {
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

  // Run target configuration - assuming there's a main.go to run
  const runTarget: TargetConfiguration = {
    executor: 'nx:run-commands',
    options: {
      command: 'go run .',
      cwd: projectRoot,
    },
  };

  // Tidy target configuration
  const tidyTarget: TargetConfiguration = {
    executor: 'nx:run-commands',
    options: {
      command: 'go mod tidy',
      cwd: projectRoot,
    },
  };

  // Lint target configuration - using golangci-lint if it's installed
  const lintTarget: TargetConfiguration = {
    executor: 'nx:run-commands',
    options: {
      command: 'golangci-lint run',
      cwd: projectRoot,
    },
  };

  // Create the project configuration
  const projectConfig: ProjectConfiguration & { root: string } = {
    name: projectName,
    root: projectRoot,
    sourceRoot: projectRoot,
    projectType: 'application',
    targets: {
      [buildTargetName]: buildTarget,
      [testTargetName]: testTarget,
      [runTargetName]: runTarget,
      [tidyTargetName]: tidyTarget,
      [lintTargetName]: lintTarget,
    },
    tags: options.tagName ? [options.tagName] : [],
  };

  return {
    projects: {
      [projectRoot]: projectConfig,
    },
  };
}
