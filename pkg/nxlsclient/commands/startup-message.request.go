package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	StartupMessageRequestMethod = "nx/startupMessage"
)

func (c *Commander) SendStartupMessageRequest(ctx context.Context, schema nxtypes.GeneratorSchema) (*nxtypes.StartupMessageDefinition, error) {
	var result *nxtypes.StartupMessageDefinition
	err := c.sendRequest(ctx, StartupMessageRequestMethod, schema, &result)
	return result, err
}
