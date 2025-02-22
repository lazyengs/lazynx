package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	ProjectsByPathsRequestMethod = "nx/projectsByPaths"
)

type ProjectsByPathsParams struct {
	Paths []string `json:"paths"`
}

func (c *Commander) SendProjectsByPathsRequest(ctx context.Context, params ProjectsByPathsParams) (map[string]*nxtypes.ProjectConfiguration, error) {
	var result map[string]*nxtypes.ProjectConfiguration
	err := c.sendRequest(ctx, ProjectsByPathsRequestMethod, params, &result)
	return result, err
}
