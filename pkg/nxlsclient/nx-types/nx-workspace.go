package nxtypes

type NxWorkspace struct {
	ProjectGraph    ProjectGraph `json:"projectGraph"`
	WorkspaceLayout struct {
		AppsDir                  *string `json:"appsDir,omitempty"`
		LibsDir                  *string `json:"libsDir,omitempty"`
		ProjectNameAndRootFormat *string `json:"projectNameAndRootFormat,omitempty"`
	} `json:"workspaceLayout"`
	SourceMaps         *ConfigurationSourceMaps `json:"sourceMaps,omitempty"`
	ProjectFileMap     *ProjectFileMap          `json:"projectFileMap,omitempty"`
	IsPartial          *bool                    `json:"isPartial,omitempty"`
	WorkspacePath      string                   `json:"workspacePath"`
	NxJson             NxJsonConfiguration      `json:"nxJson"`
	Errors             []NxError                `json:"errors,omitempty"`
	NxVersion          NxVersion                `json:"nxVersion"`
	ValidWorkspaceJson bool                     `json:"validWorkspaceJson"`
	IsLerna            bool                     `json:"isLerna"`
	IsEncapsulatedNx   bool                     `json:"isEncapsulatedNx"`
}

type NxProjectConfiguration struct {
	ProjectConfiguration
	Files []struct {
		File string `json:"file"`
	} `json:"files,omitempty"`
}

type NxError struct {
	Name    *string     `json:"name,omitempty"`
	Message *string     `json:"message,omitempty"`
	File    *string     `json:"file,omitempty"`
	Plugin  *string     `json:"plugin,omitempty"`
	Stack   *string     `json:"stack,omitempty"`
	Cause   interface{} `json:"cause,omitempty"`
}
