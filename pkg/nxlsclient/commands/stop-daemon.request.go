package commands

import (
	"context"
)

const (
	StopNxDaemonRequestMethod = "nx/stopDaemon"
)

func (c *Commander) SendStopNxDaemonRequest(ctx context.Context) error {
	var result any

	err := c.sendRequest(ctx, StopNxDaemonRequestMethod, nil, result)
	if err != nil {
		return err
	}

	return nil
}
