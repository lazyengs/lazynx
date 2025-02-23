package commands

import "context"

const (
	HasAffectedProjectsRequestMethod = "nx/hasAffectedProjects"
)

// SendHasAffectedProjectsRequest sends a request to check if there are affected projects.
func (c *Commander) SendHasAffectedProjectsRequest(ctx context.Context) (bool, error) {
	var result bool
	err := c.sendRequest(ctx, HasAffectedProjectsRequestMethod, nil, &result)
	return result, err
}
