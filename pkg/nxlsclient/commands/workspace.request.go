package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	WorkspaceRequestMethod = "nx/workspace"
)

type WorkspaceRequestParams struct {
	Reset bool `json:"reset"`
}

func (c *Commander) SendWorkspaceRequest(ctx context.Context, params *WorkspaceRequestParams) (*nxtypes.NxWorkspace, error) {
	var result *nxtypes.NxWorkspace

	err := c.sendRequest(ctx, WorkspaceRequestMethod, params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
