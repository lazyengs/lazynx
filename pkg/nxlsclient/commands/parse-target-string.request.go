package commands

import (
	"context"

	nxtypes "github.com/lazyengs/lazynx/pkg/nxlsclient/nx-types"
)

const (
	ParseTargetStringRequestMethod = "nx/parseTargetString"
)

// SendParseTargetStringRequest sends a request to parse a target string.
func (c *Commander) SendParseTargetStringRequest(ctx context.Context, targetString string) (*nxtypes.Target, error) {
	var result *nxtypes.Target
	err := c.sendRequest(ctx, ParseTargetStringRequestMethod, targetString, &result)
	return result, err
}
