package nxlsclient

import (
	"context"

	"go.lsp.dev/protocol"
)

func (c *Client) sendInitializeCommand(ctx context.Context) error {
	type InitializationOptions struct {
		workspace string
	}

	type Params struct {
		initializationOptions InitializationOptions
	}

	c.Logger.Infow("Sending initialize command")

	params := &protocol.InitializeParams{
		RootURI: protocol.DocumentURI(c.nxWorkspacePath),
		Capabilities: protocol.ClientCapabilities{
			Workspace: &protocol.WorkspaceClientCapabilities{
				Configuration: true,
			},
			TextDocument: &protocol.TextDocumentClientCapabilities{},
		},
		InitializationOptions: map[string]any{
			"workspacePath": c.nxWorkspacePath,
		},
	}

	var result any
	if err := c.conn.Call(ctx, "initialize", params, &result); err != nil {
		c.Logger.Infof("An error ocurred while executing the initialization command: %s", err.Error())
		return err
	}
	c.Logger.Infow("Initialization processed correctly", "result", result)

	return nil
}
