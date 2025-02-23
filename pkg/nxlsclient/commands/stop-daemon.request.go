package commands

import (
	"context"
)

const (
	StopNxDaemonRequestMethod = "nx/stopDaemon"
)

// SendStopNxDaemonRequest sends a request to stop the NX daemon.
func (c *Commander) SendStopNxDaemonRequest(ctx context.Context) error {
	var result any

	err := c.sendRequest(ctx, StopNxDaemonRequestMethod, nil, result)
	if err != nil {
		return err
	}

	return nil
}
