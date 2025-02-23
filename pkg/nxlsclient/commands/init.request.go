package commands

import (
	"context"

	"go.lsp.dev/protocol"
)

const (
	InitializeRequestMethod = "initialize"
)

// InitializeRequestResult represents the result of an initialize request.
type InitializeRequestResult struct {
	Capabilities struct {
		Workspace struct {
			FileOperations struct {
				DidCreate struct {
					Filters []struct {
						Pattern struct {
							Glob    string `json:"glob"`    // Glob pattern for file creation.
							Matches string `json:"matches"` // Matches specifies the type of match.
						} `json:"pattern"`
					} `json:"filters"`
				} `json:"didCreate"`
				DidDelete struct {
					Filters []struct {
						Pattern struct {
							Glob    string `json:"glob"`    // Glob pattern for file deletion.
							Matches string `json:"matches"` // Matches specifies the type of match.
						} `json:"pattern"`
					} `json:"filters"`
				} `json:"didDelete"`
			} `json:"fileOperations"`
		} `json:"workspace"`
		CompletionProvider struct {
			TriggerCharacters []string `json:"triggerCharacters"` // TriggerCharacters are characters that trigger completion.
			ResolveProvider   bool     `json:"resolveProvider"`   // ResolveProvider indicates if the server provides additional information for completion items.
		} `json:"completionProvider"`
		TextDocumentSync     int `json:"textDocumentSync"` // TextDocumentSync specifies how text documents are synced.
		DocumentLinkProvider struct {
			ResolveProvider  bool `json:"resolveProvider"`  // ResolveProvider indicates if the server provides additional information for document links.
			WorkDoneProgress bool `json:"workDoneProgress"` // WorkDoneProgress indicates if the server supports work done progress.
		} `json:"documentLinkProvider"`
		DefinitionProvider bool `json:"definitionProvider"` // DefinitionProvider indicates if the server provides definition support.
		HoverProvider      bool `json:"hoverProvider"`      // HoverProvider indicates if the server provides hover support.
	} `json:"capabilities"`
	Pid int `json:"pid"` // Pid is the process ID of the server.
}

// SendInitializeRequest sends an initialize request to the server.
func (c *Commander) SendInitializeRequest(ctx context.Context, params *protocol.InitializeParams) (*InitializeRequestResult, error) {
	var result *InitializeRequestResult
	err := c.sendRequest(ctx, InitializeRequestMethod, params, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
