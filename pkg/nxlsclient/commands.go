package nxlsclient

import (
	"context"

	"go.lsp.dev/protocol"
)

type InitializeCommandResult struct {
	Capabilities struct {
		Workspace struct {
			FileOperations struct {
				DidCreate struct {
					Filters []struct {
						Pattern struct {
							Glob    string `json:"glob"`
							Matches string `json:"matches"`
						} `json:"pattern"`
					} `json:"filters"`
				} `json:"didCreate"`
				DidDelete struct {
					Filters []struct {
						Pattern struct {
							Glob    string `json:"glob"`
							Matches string `json:"matches"`
						} `json:"pattern"`
					} `json:"filters"`
				} `json:"didDelete"`
			} `json:"fileOperations"`
		} `json:"workspace"`
		CompletionProvider struct {
			TriggerCharacters []string `json:"triggerCharacters"`
			ResolveProvider   bool     `json:"resolveProvider"`
		} `json:"completionProvider"`
		TextDocumentSync     int `json:"textDocumentSync"`
		DocumentLinkProvider struct {
			ResolveProvider  bool `json:"resolveProvider"`
			WorkDoneProgress bool `json:"workDoneProgress"`
		} `json:"documentLinkProvider"`
		DefinitionProvider bool `json:"definitionProvider"`
		HoverProvider      bool `json:"hoverProvider"`
	} `json:"capabilities"`
	Pid int `json:"pid"`
}

func (c *Client) sendInitializeCommand(ctx context.Context) (*InitializeCommandResult, error) {
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

	var result *InitializeCommandResult
	if err := c.conn.Call(ctx, "initialize", params, &result); err != nil {
		c.Logger.Infof("An error ocurred while executing the initialization command: %s", err.Error())
		return nil, err
	}
	c.Logger.Infow("Initialization processed correctly", "result", result)

	return result, nil
}
