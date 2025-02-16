package nxlsclient

import (
	"context"
	"errors"

	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
	"go.lsp.dev/protocol"
)

type WorkspaceCommandParams struct {
	protocol.WorkDoneProgressParams
}

func (c *Client) sendWorkspaceCommand(ctx context.Context, params *protocol.InitializeParams) (*nxtypes.NxWorkspace, error) {
	result, err := c.sendRequest(ctx, "nx/workspace", params)
	if err != nil {
		return nil, err
	}

	workspaceResult, ok := result.(*nxtypes.NxWorkspace)
	if !ok {
		return nil, errors.New("failed to cast result to *WorkspaceCommandResult")
	}

	return workspaceResult, nil
}
