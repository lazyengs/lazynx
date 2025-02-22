package commands

import "context"

const RefreshWorkspaceNotificationMethod = "nx/refreshWorkspace"

func (c *Commander) SendWorkspaceRefreshNotification(ctx context.Context) error {
	return c.sendNotification(ctx, RefreshWorkspaceNotificationMethod, nil)
}
