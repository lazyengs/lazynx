package commands

import (
	"context"
)

const (
	StopNxDaemonRequestMethod = "nx/stopDaemon"
)

// StopDaemonParams is the parameter structure for stop daemon request
type StopDaemonParams struct {
	Force bool `json:"force"` // Option to force stop the daemon
}

// SendStopNxDaemonRequest sends a request to stop the NX daemon.
func (c *Commander) SendStopNxDaemonRequest(ctx context.Context) error {
	var result any

	// Create parameters with force option
	params := StopDaemonParams{
		Force: true, // Force stop to ensure daemon is killed
	}

	c.Logger.Debugw("Sending stop daemon request with parameters", "params", params)
	err := c.sendRequest(ctx, StopNxDaemonRequestMethod, params, &result)
	if err != nil {
		return err
	}

	c.Logger.Debugw("Successfully stopped NX daemon via LSP")
	return nil
}
