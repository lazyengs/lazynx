package commands

import "context"

const (
	ChangeWorkspaceNotificationMethod = "nx/changeWorkspace"
)

// SendChangeWorkspaceNotification sends a notification to change the workspace.
// It takes a context and the new workspace string as parameters.
func (c *Commander) SendChangeWorkspaceNotification(ctx context.Context, workspace string) error {
	return c.sendNotification(ctx, ChangeWorkspaceNotificationMethod, workspace)
}
