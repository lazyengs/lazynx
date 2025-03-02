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

	c.Logger.Debug("Sending stop daemon request")
	err := c.sendRequest(ctx, StopNxDaemonRequestMethod, nil, &result)
	if err != nil {
		return err
	}

	c.Logger.Debugw("Successfully stopped NX daemon via LSP")
	return nil
}
