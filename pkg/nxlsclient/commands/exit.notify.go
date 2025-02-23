package commands

import (
	"context"
)

const (
	ExitNotificationMethod = "exit"
)

// Send an exit notification with the specified method and no parameters.
func (c *Commander) SendExitNotification(ctx context.Context) error {
	err := c.sendNotification(ctx, ExitNotificationMethod, nil)
	if err != nil {
		return err
	}

	return nil
}
