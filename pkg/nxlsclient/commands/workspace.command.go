package commands

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
)

const (
	WorkspaceCommandMethod = "nx/workspace"
)

type WorkspaceCommandParams struct {
	Reset bool `json:"reset"`
}

func (c *Commander) SendWorkspaceCommand(ctx context.Context, params *WorkspaceCommandParams) (*nxtypes.NxWorkspace, error) {
	var result *nxtypes.NxWorkspace

	err := c.sendRequest(ctx, WorkspaceCommandMethod, params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
