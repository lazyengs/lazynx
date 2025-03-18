package commands

import (
	"context"

	nxtypes "github.com/lazyengs/lazynx/pkg/nxlsclient/nx-types"
)

const (
	VersionRequestMethod = "nx/version"
)

// SendVersionRequest sends a request to get the NX version.
func (c *Commander) SendVersionRequest(ctx context.Context) (*nxtypes.NxVersion, error) {
	var result *nxtypes.NxVersion
	err := c.sendRequest(ctx, VersionRequestMethod, nil, &result)
	return result, err
}
