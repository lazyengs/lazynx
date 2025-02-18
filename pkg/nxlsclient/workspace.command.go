package nxlsclient

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

type WorkspaceCommandParams struct {
	Reset bool `json:"reset"`
}

func (c *Client) SendWorkspaceCommand(ctx context.Context, params *WorkspaceCommandParams) (*nxtypes.NxWorkspace, error) {
	var result *nxtypes.NxWorkspace

	err := c.sendRequest(ctx, "nx/workspace", params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
