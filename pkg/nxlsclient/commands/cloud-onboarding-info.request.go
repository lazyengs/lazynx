package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	CloudOnboardingInfoRequestMethod = "nx/cloudOnboardingInfo"
)

// SendCloudOnboardingInfoRequest sends a request to retrieve cloud onboarding information.
// It takes a context as a parameter and returns a pointer to CloudOnboardingInfo and an error.
func (c *Commander) SendCloudOnboardingInfoRequest(ctx context.Context) (*nxtypes.CloudOnboardingInfo, error) {
	var result *nxtypes.CloudOnboardingInfo
	err := c.sendRequest(ctx, CloudOnboardingInfoRequestMethod, nil, &result)
	return result, err
}
