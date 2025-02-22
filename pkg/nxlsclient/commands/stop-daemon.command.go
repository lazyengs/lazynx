package commands

import (
	"context"
)

const (
	StopNxDaemonCommandMethod = "nx/stopDaemon"
)

func (c *Commander) SendStopNxDaemonCommand(ctx context.Context) error {
	var result any

	err := c.sendRequest(ctx, StopNxDaemonCommandMethod, nil, result)
	if err != nil {
		return err
	}

	return nil
}
