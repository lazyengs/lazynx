package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	VersionRequestMethod = "nx/version"
)

func (c *Commander) SendVersionRequest(ctx context.Context) (*nxtypes.NxVersion, error) {
	var result *nxtypes.NxVersion
	err := c.sendRequest(ctx, VersionRequestMethod, nil, &result)
	return result, err
}
