package commands

import (
	"context"

	nxtypes "github.com/lazyengs/lazynx/pkg/nxlsclient/nx-types"
)

const (
	ProjectByRootRequestMethod = "nx/projectByRoot"
)

// ProjectByRootParams represents the parameters for a project by root request.
type ProjectByRootParams struct {
	ProjectRoot string `json:"projectRoot"` // ProjectRoot is the root directory of the project.
}

// SendProjectByRootRequest sends a request to get project configuration by root.
func (c *Commander) SendProjectByRootRequest(ctx context.Context, params ProjectByRootParams) (*nxtypes.ProjectConfiguration, error) {
	var result *nxtypes.ProjectConfiguration
	err := c.sendRequest(ctx, ProjectByRootRequestMethod, params, &result)
	return result, err
}
