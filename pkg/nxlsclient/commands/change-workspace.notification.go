package commands

import "context"

const (
	WorkspaceChangeNotificationMethod = "nx/changeWorkspace"
)

// SendChangeWorkspaceNotification sends a notification to change the workspace.
// It takes a context and the new workspace string as parameters.
func (c *Commander) SendChangeWorkspaceNotification(ctx context.Context, workspace string) error {
	return c.sendNotification(ctx, WorkspaceChangeNotificationMethod, workspace)
}
