package commands

import (
	"context"

	nxtypes "github.com/lazyengs/lazynx/pkg/nxlsclient/nx-types"
)

const (
	ProjectsByPathsRequestMethod = "nx/projectsByPaths"
)

// ProjectsByPathsParams represents the parameters for a projects by paths request.
type ProjectsByPathsParams struct {
	Paths []string `json:"paths"` // Paths are the paths to the projects.
}

// SendProjectsByPathsRequest sends a request to get project configurations by paths.
func (c *Commander) SendProjectsByPathsRequest(ctx context.Context, params ProjectsByPathsParams) (map[string]*nxtypes.ProjectConfiguration, error) {
	var result map[string]*nxtypes.ProjectConfiguration
	err := c.sendRequest(ctx, ProjectsByPathsRequestMethod, params, &result)
	return result, err
}
