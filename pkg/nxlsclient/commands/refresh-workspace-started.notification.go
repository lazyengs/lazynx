package commands

import "context"

const (
	WorkspaceRefreshStartedNotificationMethod = "nx/refreshWorkspaceStarted"
)

func (c *Commander) SendWorkspaceRefreshStartedNotification(ctx context.Context) error {
	return c.sendNotification(ctx, WorkspaceRefreshStartedNotificationMethod, nil)
}
