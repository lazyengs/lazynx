package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	ProjectByRootRequestMethod = "nx/projectByRoot"
)

type ProjectByRootParams struct {
	ProjectRoot string `json:"projectRoot"`
}

func (c *Commander) SendProjectByRootRequest(ctx context.Context, params ProjectByRootParams) (*nxtypes.ProjectConfiguration, error) {
	var result *nxtypes.ProjectConfiguration
	err := c.sendRequest(ctx, ProjectByRootRequestMethod, params, &result)
	return result, err
}
