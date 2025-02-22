package commands

import (
	"context"

	"go.lsp.dev/protocol"
)

const (
	InitializeCommandMethod = "initialize"
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

func (c *Commander) SendInitializeCommand(ctx context.Context, params *protocol.InitializeParams) (*InitializeCommandResult, error) {
	var result *InitializeCommandResult
	err := c.sendRequest(ctx, InitializeCommandMethod, params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
