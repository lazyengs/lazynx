package nxtypes

type NxJsonConfiguration struct {
	DefaultProject  *string                      `json:"defaultProject,omitempty"`
	NamedInputs     map[string][]interface{}     `json:"namedInputs,omitempty"`
	Installation    *NxInstallationConfiguration `json:"installation,omitempty"`
	Release         *NxReleaseConfiguration      `json:"release,omitempty"`
	Affected        *NxAffectedConfig            `json:"affected,omitempty"`
	DefaultBase     *string                      `json:"defaultBase,omitempty"`
	WorkspaceLayout *struct {
		LibsDir *string `json:"libsDir,omitempty"`
		AppsDir *string `json:"appsDir,omitempty"`
	} `json:"workspaceLayout,omitempty"`
	TasksRunnerOptions map[string]struct {
		Runner  *string     `json:"runner,omitempty"`
		Options interface{} `json:"options,omitempty"`
	} `json:"tasksRunnerOptions,omitempty"`
	Generators map[string]map[string]interface{} `json:"generators,omitempty"`
	Cli        *struct {
		PackageManager     *PackageManager `json:"packageManager,omitempty"`
		DefaultProjectName *string         `json:"defaultProjectName,omitempty"`
	} `json:"cli,omitempty"`
	Extends              *string                           `json:"extends,omitempty"`
	PluginsConfig        map[string]map[string]interface{} `json:"pluginsConfig,omitempty"`
	UseLegacyCache       *bool                             `json:"useLegacyCache,omitempty"`
	ImplicitDependencies map[string]interface{}            `json:"implicitDependencies,omitempty"`
	TargetDefaults       TargetDefaults                    `json:"targetDefaults,omitempty"`
	NxCloudAccessToken   *string                           `json:"nxCloudAccessToken,omitempty"`
	NxCloudId            *string                           `json:"nxCloudId,omitempty"`
	NxCloudUrl           *string                           `json:"nxCloudUrl,omitempty"`
	NxCloudEncryptionKey *string                           `json:"nxCloudEncryptionKey,omitempty"`
	Parallel             *int                              `json:"parallel,omitempty"`
	CacheDirectory       *string                           `json:"cacheDirectory,omitempty"`
	UseDaemonProcess     *bool                             `json:"useDaemonProcess,omitempty"`
	UseInferencePlugins  *bool                             `json:"useInferencePlugins,omitempty"`
	NeverConnectToCloud  *bool                             `json:"neverConnectToCloud,omitempty"`
	Sync                 *NxSyncConfiguration              `json:"sync,omitempty"`
	Plugins              []PluginConfiguration             `json:"plugins,omitempty"`
}

type NxInstallationConfiguration struct {
	Plugins map[string]string `json:"plugins,omitempty"`
	Version string            `json:"version"`
}

type NxReleaseChangelogConfiguration struct {
	CreateRelease      interface{}             `json:"createRelease,omitempty"`      // false or "github" or struct
	EntryWhenNoChanges interface{}             `json:"entryWhenNoChanges,omitempty"` // string or false
	File               interface{}             `json:"file,omitempty"`               // string or false
	Renderer           *string                 `json:"renderer,omitempty"`
	RenderOptions      *ChangelogRenderOptions `json:"renderOptions,omitempty"`
}

type NxReleaseConfiguration struct {
	Projects interface{} `json:"projects,omitempty"` // []string or string
	Groups   map[string]struct {
		ProjectsRelationship *string                        `json:"projectsRelationship,omitempty"`
		Projects             interface{}                    `json:"projects"` // []string or string
		Version              *NxReleaseVersionConfiguration `json:"version,omitempty"`
		Changelog            interface{}                    `json:"changelog,omitempty"` // bool or NxReleaseChangelogConfiguration
		ReleaseTagPattern    *string                        `json:"releaseTagPattern,omitempty"`
		VersionPlans         interface{}                    `json:"versionPlans,omitempty"` // bool or NxReleaseVersionPlansConfiguration
	} `json:"groups,omitempty"`
	ProjectsRelationship *string `json:"projectsRelationship,omitempty"`
	Changelog            *struct {
		Git                *NxReleaseGitConfiguration `json:"git,omitempty"`
		WorkspaceChangelog interface{}                `json:"workspaceChangelog,omitempty"`
		ProjectChangelogs  interface{}                `json:"projectChangelogs,omitempty"`
		AutomaticFromRef   *bool                      `json:"automaticFromRef,omitempty"`
	} `json:"changelog,omitempty"`
	Version             *NxReleaseVersionConfiguration             `json:"version,omitempty"`
	ReleaseTagPattern   *string                                    `json:"releaseTagPattern,omitempty"`
	Git                 *NxReleaseGitConfiguration                 `json:"git,omitempty"`
	ConventionalCommits *NxReleaseConventionalCommitsConfiguration `json:"conventionalCommits,omitempty"`
	VersionPlans        interface{}                                `json:"versionPlans,omitempty"`
}

type ChangelogRenderOptions map[string]interface{}

type NxSyncConfiguration struct {
	GlobalGenerators           []string                          `json:"globalGenerators,omitempty"`
	GeneratorOptions           map[string]map[string]interface{} `json:"generatorOptions,omitempty"`
	ApplyChanges               *bool                             `json:"applyChanges,omitempty"`
	DisabledTaskSyncGenerators []string                          `json:"disabledTaskSyncGenerators,omitempty"`
}

type NxReleaseVersionPlansConfiguration struct {
	IgnorePatternsForPlanCheck []string `json:"ignorePatternsForPlanCheck,omitempty"`
}

type NxReleaseGitConfiguration struct {
	Commit        *bool       `json:"commit,omitempty"`
	CommitMessage *string     `json:"commitMessage,omitempty"`
	CommitArgs    interface{} `json:"commitArgs,omitempty"` // string or []string
	StageChanges  *bool       `json:"stageChanges,omitempty"`
	Tag           *bool       `json:"tag,omitempty"`
	TagMessage    *string     `json:"tagMessage,omitempty"`
	TagArgs       interface{} `json:"tagArgs,omitempty"` // string or []string
	Push          *bool       `json:"push,omitempty"`
}

type NxReleaseConventionalCommitsConfiguration struct {
	Types map[string]interface{} `json:"types,omitempty"` // value can be bool or struct
}

type ImplicitJsonSubsetDependency map[string]interface{}

type NxReleaseVersionConfiguration struct {
	Generator           *string                `json:"generator,omitempty"`
	GeneratorOptions    map[string]interface{} `json:"generatorOptions,omitempty"`
	ConventionalCommits *bool                  `json:"conventionalCommits,omitempty"`
}

type NxAffectedConfig struct {
	DefaultBase *string `json:"defaultBase,omitempty"`
}

type TargetDefaults map[string]TargetConfiguration

type PluginConfiguration interface{} // string or ExpandedPluginConfiguration

type ExpandedPluginConfiguration struct {
	Plugin  string      `json:"plugin"`
	Options interface{} `json:"options,omitempty"`
	Include []string    `json:"include,omitempty"`
	Exclude []string    `json:"exclude,omitempty"`
}
