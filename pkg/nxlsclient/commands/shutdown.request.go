package commands

import (
	"context"
)

const (
	ShutdownRequestMethod = "shutdown"
)

// SendShutdownRequest sends a shutdown request to the nxls server.
func (c *Commander) SendShutdownRequest(ctx context.Context) error {
	var result any

	// For LSP shutdown we don't need parameters
	c.Logger.Debugw("Sending shutdown request to LSP server")
	err := c.sendRequest(ctx, ShutdownRequestMethod, nil, &result)
	if err != nil {
		return err
	}

	c.Logger.Debugw("LSP server shutdown successful")
	return nil
}
