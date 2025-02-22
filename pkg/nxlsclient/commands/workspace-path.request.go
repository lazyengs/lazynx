package commands

import "context"

const (
	WorkspacePathRequestMethod = "nx/workspacePath"
)

func (c *Commander) SendWorkspacePathRequest(ctx context.Context) (string, error) {
	var result string
	err := c.sendRequest(ctx, WorkspacePathRequestMethod, nil, &result)
	return result, err
}
