package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	ParseTargetStringRequestMethod = "nx/parseTargetString"
)

func (c *Commander) SendParseTargetStringRequest(ctx context.Context, targetString string) (*nxtypes.Target, error) {
	var result *nxtypes.Target
	err := c.sendRequest(ctx, ParseTargetStringRequestMethod, targetString, &result)
	return result, err
}
