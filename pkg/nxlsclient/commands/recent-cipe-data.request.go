package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	RecentCIPEDataRequestMethod = "nx/recentCIPEData"
)

type RecentCIPEDataResult struct {
	Info         []nxtypes.CIPEInfo     `json:"info,omitempty"`
	Error        *nxtypes.CIPEInfoError `json:"error,omitempty"`
	WorkspaceUrl string                 `json:"workspaceUrl,omitempty"`
}

func (c *Commander) SendRecentCIPEDataRequest(ctx context.Context) (*RecentCIPEDataResult, error) {
	var result *RecentCIPEDataResult
	err := c.sendRequest(ctx, RecentCIPEDataRequestMethod, nil, &result)
	return result, err
}
