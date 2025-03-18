package commands

import (
	"context"

	nxtypes "github.com/lazyengs/lazynx/pkg/nxlsclient/nx-types"
)

const (
	StartupMessageRequestMethod = "nx/startupMessage"
)

// SendStartupMessageRequest sends a request to get the startup message.
func (c *Commander) SendStartupMessageRequest(ctx context.Context, schema nxtypes.GeneratorSchema) (*nxtypes.StartupMessageDefinition, error) {
	var result *nxtypes.StartupMessageDefinition
	err := c.sendRequest(ctx, StartupMessageRequestMethod, schema, &result)
	return result, err
}
