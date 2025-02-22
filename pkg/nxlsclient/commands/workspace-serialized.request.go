package commands

import "context"

const (
	WorkspaceSerializedCommandMethod = "nx/workspaceSerialized"
)

func (c *Commander) SendWorkspaceSerializedCommand(ctx context.Context, params WorkspaceCommandParams) (string, error) {
	var result string
	err := c.sendRequest(ctx, WorkspaceSerializedCommandMethod, params, &result)
	return result, err
}
