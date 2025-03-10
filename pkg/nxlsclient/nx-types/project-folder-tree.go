package nxtypes

type TreeNode struct {
	Dir                  string                   `json:"dir"`
	ProjectName          string                   `json:"projectName,omitempty"`
	ProjectConfiguration *ProjectGraphProjectNode `json:"projectConfiguration,omitempty"`
	Children             []string                 `json:"children"`
}

// TreeMap is a map of string to TreeNode
type TreeMap map[string]TreeNode
