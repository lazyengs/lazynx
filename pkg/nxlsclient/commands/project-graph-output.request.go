package commands

import (
	"context"
)

const (
	ProjectGraphOutputRequestMethod = "nx/projectGraphOutput"
)

// ProjectGraphOutputResult represents the result of a project graph output request.
type ProjectGraphOutputResult struct {
	Directory    string `json:"directory"`    // Directory is the directory of the project graph output.
	RelativePath string `json:"relativePath"` // RelativePath is the relative path of the project graph output.
	FullPath     string `json:"fullPath"`     // FullPath is the full path of the project graph output.
}

// SendProjectGraphOutputRequest sends a request to get the project graph output.
func (c *Commander) SendProjectGraphOutputRequest(ctx context.Context) (*ProjectGraphOutputResult, error) {
	var result *ProjectGraphOutputResult
	err := c.sendRequest(ctx, ProjectGraphOutputRequestMethod, nil, &result)
	return result, err
}
