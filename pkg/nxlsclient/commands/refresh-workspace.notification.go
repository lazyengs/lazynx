package commands

import "context"

const RefreshWorkspaceNotificationMethod = "nx/refreshWorkspace"

// SendWorkspaceRefreshNotification sends a notification to refresh the workspace.
func (c *Commander) SendWorkspaceRefreshNotification(ctx context.Context) error {
	return c.sendNotification(ctx, RefreshWorkspaceNotificationMethod, nil)
}
