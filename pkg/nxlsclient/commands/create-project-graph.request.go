package commands

import "context"

const (
	CreateProjectGraphRequestMethod = "nx/createProjectGraph"
)

type CreateProjectGraphParams struct {
	ShowAffected bool `json:"showAffected"`
}

func (c *Commander) SendCreateProjectGraphRequest(ctx context.Context, params CreateProjectGraphParams) (*string, error) {
	var result *string
	err := c.sendRequest(ctx, CreateProjectGraphRequestMethod, params, &result)
	return result, err
}
