package nxtypes

type ProjectConfiguration struct {
	Name        *string                           `json:"name,omitempty"`
	Targets     map[string]TargetConfiguration    `json:"targets,omitempty"`
	SourceRoot  *string                           `json:"sourceRoot,omitempty"`
	ProjectType *ProjectType                      `json:"projectType,omitempty"`
	Generators  map[string]map[string]interface{} `json:"generators,omitempty"`
	NamedInputs map[string][]interface{}          `json:"namedInputs,omitempty"`
	Release     *struct {
		Version *struct {
			Generator        *string                `json:"generator,omitempty"`
			GeneratorOptions map[string]interface{} `json:"generatorOptions,omitempty"`
		} `json:"version,omitempty"`
	} `json:"release,omitempty"`
	Metadata             *ProjectMetadata `json:"metadata,omitempty"`
	Root                 string           `json:"root"`
	ImplicitDependencies []string         `json:"implicitDependencies,omitempty"`
	Tags                 []string         `json:"tags,omitempty"`
}

type InputDefinition struct {
	Projects                  interface{} `json:"projects,omitempty"`
	Input                     *string     `json:"input,omitempty"`
	Dependencies              *bool       `json:"dependencies,omitempty"`
	Fileset                   *string     `json:"fileset,omitempty"`
	Runtime                   *string     `json:"runtime,omitempty"`
	DependentTasksOutputFiles *string     `json:"dependentTasksOutputFiles,omitempty"`
	Transitive                *bool       `json:"transitive,omitempty"`
	Env                       *string     `json:"env,omitempty"`
	ExternalDependencies      []string    `json:"externalDependencies,omitempty"`
}

type TargetConfiguration struct {
	Executor             *string                `json:"executor,omitempty"`
	Command              *string                `json:"command,omitempty"`
	Outputs              []string               `json:"outputs,omitempty"`
	DependsOn            []interface{}          `json:"dependsOn,omitempty"` // []TargetDependencyConfig or []string
	Inputs               []interface{}          `json:"inputs,omitempty"`    // []InputDefinition or []string
	Options              interface{}            `json:"options,omitempty"`
	Configurations       map[string]interface{} `json:"configurations,omitempty"`
	DefaultConfiguration *string                `json:"defaultConfiguration,omitempty"`
	Cache                *bool                  `json:"cache,omitempty"`
	Metadata             *TargetMetadata        `json:"metadata,omitempty"`
	Parallelism          *bool                  `json:"parallelism,omitempty"`
	SyncGenerators       []string               `json:"syncGenerators,omitempty"`
}

type TargetMetadata struct {
	Description       *string `json:"description,omitempty"`
	NonAtomizedTarget *string `json:"nonAtomizedTarget,omitempty"`
	Help              *struct {
		Command string `json:"command"`
		Example struct {
			Options map[string]interface{} `json:"options,omitempty"`
			Args    []string               `json:"args,omitempty"`
		} `json:"example"`
	} `json:"help,omitempty"`
	Technologies []string `json:"technologies,omitempty"`
}

type TargetDependencyConfig struct {
	Projects     interface{} `json:"projects,omitempty"`
	Dependencies *bool       `json:"dependencies,omitempty"`
	Params       *string     `json:"params,omitempty"`
	Target       string      `json:"target"`
}

type ProjectMetadata struct {
	Description  *string             `json:"description,omitempty"`
	TargetGroups map[string][]string `json:"targetGroups,omitempty"`
	Owners       map[string]struct {
		OwnedFiles []struct {
			Files      interface{} `json:"files"`
			FromConfig *struct {
				FilePath string `json:"filePath"`
				Location struct {
					StartLine int `json:"startLine"`
					EndLine   int `json:"endLine"`
				} `json:"location"`
			} `json:"fromConfig,omitempty"`
		} `json:"ownedFiles"`
	} `json:"owners,omitempty"`
	Js *struct {
		PackageExports interface{} `json:"packageExports,omitempty"`
		PackageName    string      `json:"packageName"`
	} `json:"js,omitempty"`
	Technologies []string `json:"technologies,omitempty"` // ["*"] or []string
}

type ProjectType string

const (
	ProjectTypeLibrary     ProjectType = "library"
	ProjectTypeApplication ProjectType = "application"
)
