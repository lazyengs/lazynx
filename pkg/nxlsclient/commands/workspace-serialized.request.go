package commands

import "context"

const (
	WorkspaceSerializedRequestMethod = "nx/workspaceSerialized"
)

func (c *Commander) SendWorkspaceSerializedRequest(ctx context.Context, params WorkspaceRequestParams) (string, error) {
	var result string
	err := c.sendRequest(ctx, WorkspaceSerializedRequestMethod, params, &result)
	return result, err
}
