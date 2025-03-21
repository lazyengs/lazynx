import {
  CreateNodesContextV2,
  CreateNodesV2,
  ProjectConfiguration,
  TargetConfiguration,
  createNodesFromFiles,
  joinPathFragments,
} from '@nx/devkit';
import { basename } from 'path';

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
      command: 'go build ./...',
      cwd: projectRoot,
    },
    cache: true,
    inputs: [
      '{projectRoot}/go.mod',
      '{projectRoot}/go.sum',
      joinPathFragments('{projectRoot}', '**', '*.go'),
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

import {
  CreateDependencies,
  CreateDependenciesContext,
  DependencyType,
  FileData,
  RawProjectGraphDependency,
  workspaceRoot,
} from '@nx/devkit';
import { execSync } from 'child_process';
import { readFileSync } from 'fs';
import { dirname, extname } from 'path';

export interface GoPluginOptions {
  buildTargetName?: string;
  testTargetName?: string;
  runTargetName?: string;
  tidyTargetName?: string;
  lintTargetName?: string;
  tagName?: string;
}

export type GoListType = 'import' | 'use';

const REGEXS: Record<GoListType | 'version', RegExp> = {
  import: /import\s+(?:(\w+)\s+)?"([^"]+)"|\(([\s\S]*?)\)/,
  use: /use\s+(\(([^)]*)\)|([^\n]*))/,
  version: /go(?<version>\S+) /,
};

/**
 * Executes the `go list -m -json` command in the
 * specified directory and returns the output as a string.
 *
 * @param cwd the current working directory where the command should be executed.
 * @param failSilently if true, the function will return an empty string instead of throwing an error when the command fails.
 * @returns The output of the `go list -m -json` command as a string.
 * @throws Will throw an error if the command fails and `failSilently` is false.
 */
export const getGoModules = (cwd: string, failSilently: boolean): string => {
  try {
    return execSync('go list -m -json', {
      encoding: 'utf-8',
      cwd,
      stdio: ['ignore'],
      windowsHide: true,
    });
  } catch (error) {
    if (failSilently) {
      return '';
    } else {
      throw error;
    }
  }
};

/**
 * Parses a Go list (also support list with only one item).
 *
 * @param listType type of list to parse
 * @param content list to parse as a string
 */
export const parseGoList = (
  listType: GoListType,
  content: string
): string[] => {
  const exec = REGEXS[listType].exec(content);
  return (
    (exec?.[2] ?? exec?.[3])
      ?.trim()
      .split(/\n+/)
      .map((line) => line.trim()) ?? []
  );
};

type ProjectRootMap = Map<string, string>;

interface GoModule {
  Path: string;
  Dir: string;
}

interface GoImportWithModule {
  import: string;
  module: GoModule;
}

/**
 * Computes a list of go modules.
 *
 * @param failSilently if true, the function will not throw an error if it fails
 */
const computeGoModules = (failSilently = false): GoModule[] => {
  const blocks = getGoModules(workspaceRoot, failSilently);
  if (blocks != null) {
    return blocks
      .split('}')
      .filter((block) => block.trim().length > 0)
      .map((block) => JSON.parse(`${block}}`) as GoModule)
      .sort((module1, module2) => module1.Path.localeCompare(module2.Path))
      .reverse();
  }
  throw new Error('Cannot get list of Go modules');
};

/**
 * Extracts a map of project root to project name based on context.
 *
 * @param context the Nx graph context
 */
const extractProjectRootMap = (
  context: CreateDependenciesContext
): ProjectRootMap =>
  Object.keys(context.projects).reduce((map, name) => {
    map.set(context.projects[name].root, name);
    return map;
  }, new Map<string, string>());

/**
 * Gets a list of go imports with associated module in the file.
 *
 * @param fileData file object computed by Nx
 * @param modules list of go modules
 */
const getFileModuleImports = (
  fileData: FileData,
  modules: GoModule[]
): GoImportWithModule[] => {
  const content = readFileSync(fileData.file, 'utf-8')?.toString();
  if (content == null) {
    return [];
  }
  return parseGoList('import', content)
    .map((item) => (item.includes('"') ? item.split('"')[1] : item))
    .filter((item) => item != null)
    .map((item) => ({
      import: item,
      module: modules.find((mod) => item.startsWith(mod.Path)),
    }))
    .filter((item) => item.module);
};

/**
 * Gets the project name for the go import by getting the relative path for the import with in the go module system
 * then uses that to calculate the relative path on disk and looks up which project in the workspace the import is a part
 * of.
 *
 * @param projectRootMap map with project roots in the workspace
 * @param import the go import
 * @param module the go module
 */
const getProjectNameForGoImport = (
  projectRootMap: ProjectRootMap,
  { import: goImport, module }: GoImportWithModule
): string | null => {
  const relativeImportPath = goImport.substring(module.Path.length + 1);
  const relativeModuleDir = module.Dir.substring(
    workspaceRoot.length + 1
  ).replace(/\\/g, '/');
  let projectPath = relativeModuleDir
    ? `${relativeModuleDir}/${relativeImportPath}`
    : relativeImportPath;

  while (projectPath !== '.') {
    if (projectPath.endsWith('/')) {
      projectPath = projectPath.slice(0, -1);
    }

    const projectName = projectRootMap.get(projectPath);
    if (projectName) {
      return projectName;
    }
    projectPath = dirname(projectPath);
  }
  return null;
};

export const createDependencies: CreateDependencies<GoPluginOptions> = async (
  _,
  context
) => {
  const dependencies: RawProjectGraphDependency[] = [];

  let goModules: GoModule[] = null;
  let projectRootMap: ProjectRootMap = null;

  for (const projectName in context.filesToProcess.projectFileMap) {
    const files = context.filesToProcess.projectFileMap[projectName].filter(
      (file) => extname(file.file) === '.go'
    );

    if (files.length > 0 && goModules == null) {
      goModules = computeGoModules();
      projectRootMap = extractProjectRootMap(context);
    }

    for (const file of files) {
      dependencies.push(
        ...getFileModuleImports(file, goModules)
          .map((goImport) =>
            getProjectNameForGoImport(projectRootMap, goImport)
          )
          .filter((target) => target != null)
          .map((target) => ({
            type: DependencyType.static,
            source: projectName,
            target: target,
            sourceFile: file.file,
          }))
      );
    }
  }
  return dependencies;
};
