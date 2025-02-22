package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	CloudOnboardingInfoRequestMethod = "nx/cloudOnboardingInfo"
)

func (c *Commander) SendCloudOnboardingInfoRequest(ctx context.Context) (*nxtypes.CloudOnboardingInfo, error) {
	var result *nxtypes.CloudOnboardingInfo
	err := c.sendRequest(ctx, CloudOnboardingInfoRequestMethod, nil, &result)
	return result, err
}
