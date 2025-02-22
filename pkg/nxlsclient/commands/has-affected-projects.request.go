package commands

import "context"

const (
	HasAffectedProjectsRequestMethod = "nx/hasAffectedProjects"
)

func (c *Commander) SendHasAffectedProjectsRequest(ctx context.Context) (bool, error) {
	var result bool
	err := c.sendRequest(ctx, HasAffectedProjectsRequestMethod, nil, &result)
	return result, err
}
