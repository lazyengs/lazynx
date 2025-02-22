package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	ProjectByPathRequestMethod = "nx/projectByPath"
)

type ProjectByPathParams struct {
	ProjectPath string `json:"projectPath"`
}

func (c *Commander) SendProjectByPathRequest(ctx context.Context, params ProjectByPathParams) (*nxtypes.ProjectConfiguration, error) {
	var result *nxtypes.ProjectConfiguration
	err := c.sendRequest(ctx, ProjectByPathRequestMethod, params, &result)
	return result, err
}
