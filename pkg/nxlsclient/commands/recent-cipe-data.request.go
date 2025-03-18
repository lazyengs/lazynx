package commands

import (
	"context"

	nxtypes "github.com/lazyengs/lazynx/pkg/nxlsclient/nx-types"
)

const (
	RecentCIPEDataRequestMethod = "nx/recentCIPEData"
)

// RecentCIPEDataResult represents the result of a recent CIPE data request.
type RecentCIPEDataResult struct {
	Info         []nxtypes.CIPEInfo     `json:"info,omitempty"`         // Info contains the CIPE information.
	Error        *nxtypes.CIPEInfoError `json:"error,omitempty"`        // Error contains the CIPE error information.
	WorkspaceUrl string                 `json:"workspaceUrl,omitempty"` // WorkspaceUrl is the URL of the workspace.
}

// SendRecentCIPEDataRequest sends a request to get recent CIPE data.
func (c *Commander) SendRecentCIPEDataRequest(ctx context.Context) (*RecentCIPEDataResult, error) {
	var result *RecentCIPEDataResult
	err := c.sendRequest(ctx, RecentCIPEDataRequestMethod, nil, &result)
	return result, err
}
