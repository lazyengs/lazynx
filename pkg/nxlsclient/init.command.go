package nxlsclient

import (
	"context"
	"fmt"

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

func (c *Client) sendInitializeCommand(ctx context.Context, params *protocol.InitializeParams) (*InitializeCommandResult, error) {
	result, err := c.sendRequest(ctx, "initialize", params)
	if err != nil {
		return nil, err
	}

	initializeResult, ok := result.(*InitializeCommandResult)
	if !ok {
		return nil, fmt.Errorf("failed to cast result to *InitializeCommandResult: %w", err)
	}

	return initializeResult, nil
}
