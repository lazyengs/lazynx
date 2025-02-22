package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	GeneratorContextFromPathRequestMethod = "nx/generatorContextFromPath"
)

type GeneratorContextFromPathParams struct {
	Generator *nxtypes.TaskExecutionSchema `json:"generator,omitempty"`
	Path      string                       `json:"path"`
}

type GeneratorContextFromPathResult struct {
	Path        string `json:"path,omitempty"`
	Directory   string `json:"directory,omitempty"`
	Project     string `json:"project,omitempty"`
	ProjectName string `json:"projectName,omitempty"`
}

func (c *Commander) SendGeneratorContextFromPathRequest(ctx context.Context, params GeneratorContextFromPathParams) (*GeneratorContextFromPathResult, error) {
	var result *GeneratorContextFromPathResult
	err := c.sendRequest(ctx, GeneratorContextFromPathRequestMethod, params, &result)
	return result, err
}
