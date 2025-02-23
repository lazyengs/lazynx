package commands

import "context"

const (
	WorkspaceRefreshStartedNotificationMethod = "nx/refreshWorkspaceStarted"
)

// SendWorkspaceRefreshStartedNotification sends a notification that the workspace refresh has started.
func (c *Commander) SendWorkspaceRefreshStartedNotification(ctx context.Context) error {
	return c.sendNotification(ctx, WorkspaceRefreshStartedNotificationMethod, nil)
}
