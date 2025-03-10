package nxtypes

type FileDataDependency interface{} // string or [2]string or [3]string

type FileData struct {
	File string               `json:"file"`
	Hash string               `json:"hash"`
	Deps []FileDataDependency `json:"deps,omitempty"`
}

type DependencyType string

const (
	DependencyTypeStatic   DependencyType = "static"
	DependencyTypeDynamic  DependencyType = "dynamic"
	DependencyTypeImplicit DependencyType = "implicit"
)

type ProjectGraphDependency struct {
	Type   string `json:"type"`
	Target string `json:"target"`
	Source string `json:"source"`
}

type ProjectGraphExternalNode struct {
	Data struct {
		Hash        *string `json:"hash,omitempty"`
		Version     string  `json:"version"`
		PackageName string  `json:"packageName"`
	} `json:"data"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type ProjectGraphProjectNode struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Data struct {
		Description *string `json:"description,omitempty"`
		ProjectConfiguration
	} `json:"data"`
}

type ProjectGraph struct {
	Nodes         map[string]ProjectGraphProjectNode  `json:"nodes"`
	ExternalNodes map[string]ProjectGraphExternalNode `json:"externalNodes,omitempty"`
	Dependencies  map[string][]ProjectGraphDependency `json:"dependencies"`
	Version       *string                             `json:"version,omitempty"`
}

type ProjectFileMap map[string][]FileData
