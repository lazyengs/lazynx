package commands

import "context"

const (
	WorkspaceSerializedRequestMethod = "nx/workspaceSerialized"
)

// SendWorkspaceSerializedRequest sends a request to get the serialized workspace.
func (c *Commander) SendWorkspaceSerializedRequest(ctx context.Context, params WorkspaceRequestParams) (string, error) {
	var result string
	err := c.sendRequest(ctx, WorkspaceSerializedRequestMethod, params, &result)
	return result, err
}
