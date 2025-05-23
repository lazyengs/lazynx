package commands

import (
	"context"

	nxtypes "github.com/lazyengs/lazynx/pkg/nxlsclient/nx-types"
)

const (
	ProjectByPathRequestMethod = "nx/projectByPath"
)

// ProjectByPathParams represents the parameters for a project by path request.
type ProjectByPathParams struct {
	ProjectPath string `json:"projectPath"` // ProjectPath is the path to the project.
}

// SendProjectByPathRequest sends a request to get project configuration by path.
func (c *Commander) SendProjectByPathRequest(ctx context.Context, params ProjectByPathParams) (*nxtypes.ProjectConfiguration, error) {
	var result *nxtypes.ProjectConfiguration
	err := c.sendRequest(ctx, ProjectByPathRequestMethod, params, &result)
	return result, err
}
