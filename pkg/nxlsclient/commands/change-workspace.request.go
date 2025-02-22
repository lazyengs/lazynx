package commands

import "context"

const (
	WorkspaceChangeNotificationMethod = "nx/changeWorkspace"
)

func (c *Commander) SendChangeWorkspaceNotification(ctx context.Context, workspace string) error {
	return c.sendNotification(ctx, WorkspaceChangeNotificationMethod, workspace)
}
