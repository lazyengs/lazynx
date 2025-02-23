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

	err := c.sendRequest(ctx, ShutdownRequestMethod, nil, result)
	if err != nil {
		return err
	}

	return nil
}
