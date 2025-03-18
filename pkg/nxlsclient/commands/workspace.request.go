package commands

import (
	"context"

	nxtypes "github.com/lazyengs/lazynx/pkg/nxlsclient/nx-types"
)

const (
	WorkspaceRequestMethod = "nx/workspace"
)

// WorkspaceRequestParams represents the parameters for a workspace request.
type WorkspaceRequestParams struct {
	Reset bool `json:"reset"` // Reset specifies whether to reset the workspace.
}

// SendWorkspaceRequest sends a request to get the workspace.
func (c *Commander) SendWorkspaceRequest(ctx context.Context, params *WorkspaceRequestParams) (*nxtypes.NxWorkspace, error) {
	var result *nxtypes.NxWorkspace

	err := c.sendRequest(ctx, WorkspaceRequestMethod, params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
