package nxlsclient

import (
	"context"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
	"go.lsp.dev/protocol"
)

type WorkspaceCommandParams struct {
	protocol.WorkDoneProgressParams
}

func (c *Client) sendWorkspaceCommand(ctx context.Context, params *protocol.InitializeParams) (*nxtypes.NxWorkspace, error) {
	var result *nxtypes.NxWorkspace

	err := c.sendRequest(ctx, "nx/workspace", params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
