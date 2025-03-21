import {
  CreateNodesContextV2,
  CreateNodesV2,
  TargetConfiguration,
  createNodesFromFiles,
  joinPathFragments,
} from '@nx/devkit';
import { dirname } from 'path';

// Expected format of the plugin options defined in nx.json
export interface GoPluginOptions {
  buildTargetName?: string;
  testTargetName?: string;
  runTargetName?: string;
  tidyTargetName?: string;
  lintTargetName?: string;
}

// File glob to find all Go projects
const goModGlob = '**/go.mod';

// Entry function that Nx calls to modify the graph
export const createNodesV2: CreateNodesV2<GoPluginOptions> = [
  goModGlob,
  async (configFiles, options, context) => {
    return await createNodesFromFiles(
      (configFile, options, context) =>
        createNodesInternal(configFile, options, context),
      configFiles,
      options,
      context
    );
  },
];

async function createNodesInternal(
  configFilePath: string,
  options: GoPluginOptions,
  context: CreateNodesContextV2
) {
  const projectRoot = dirname(configFilePath);

  // For better UX, set default target names if not provided
  const buildTargetName = options.buildTargetName || 'build';
  const testTargetName = options.testTargetName || 'test';
  const runTargetName = options.runTargetName || 'serve';
  const tidyTargetName = options.tidyTargetName || 'tidy';
  const lintTargetName = options.lintTargetName || 'lint';

  // Build target configuration
  const buildTarget: TargetConfiguration = {
    command: 'go build ./...',
    options: { cwd: projectRoot },
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

  // Test target configuration
  const testTarget: TargetConfiguration = {
    command: 'go test ./...',
    options: { cwd: projectRoot },
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
    command: 'go run .',
    options: { cwd: projectRoot },
  };

  // Tidy target configuration
  const tidyTarget: TargetConfiguration = {
    command: 'go mod tidy',
    options: { cwd: projectRoot },
  };

  // Lint target configuration - using golangci-lint if it's installed
  const lintTarget: TargetConfiguration = {
    command: 'golangci-lint run',
    options: { cwd: projectRoot },
  };

  // Project configuration to be merged into the rest of the Nx configuration
  return {
    projects: {
      [projectRoot]: {
        targets: {
          [buildTargetName]: buildTarget,
          [testTargetName]: testTarget,
          [runTargetName]: runTarget,
          [tidyTargetName]: tidyTarget,
          [lintTargetName]: lintTarget,
        },
      },
    },
  };
}

