package commands

import (
	"context"
)

const (
	ProjectGraphOutputRequestMethod = "nx/projectGraphOutput"
)

type ProjectGraphOutputResult struct {
	Directory    string `json:"directory"`
	RelativePath string `json:"relativePath"`
	FullPath     string `json:"fullPath"`
}

func (c *Commander) SendProjectGraphOutputRequest(ctx context.Context) (*ProjectGraphOutputResult, error) {
	var result *ProjectGraphOutputResult
	err := c.sendRequest(ctx, ProjectGraphOutputRequestMethod, nil, &result)
	return result, err
}
