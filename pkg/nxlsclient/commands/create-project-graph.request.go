package commands

import "context"

const (
	CreateProjectGraphRequestMethod = "nx/createProjectGraph"
)

// CreateProjectGraphParams represents the parameters for creating a project graph.
type CreateProjectGraphParams struct {
	ShowAffected bool `json:"showAffected"` // ShowAffected specifies whether to show affected projects.
}

// SendCreateProjectGraphRequest sends a request to create a project graph.
// It takes a context and CreateProjectGraphParams as parameters and returns a pointer to a string and an error.
func (c *Commander) SendCreateProjectGraphRequest(ctx context.Context, params CreateProjectGraphParams) (*string, error) {
	var result *string
	err := c.sendRequest(ctx, CreateProjectGraphRequestMethod, params, &result)
	return result, err
}
