package nxlsclient

import (
	"context"
	"encoding/json"

	"github.com/sourcegraph/jsonrpc2"
)

type windowLogMessageNotification struct {
	Message string `json:"message"`
	Type    int8   `json:"type"`
}

// connectToLSPServer connects to the LSP server using the provided ReadWriteCloser.
func (c *Client) connectToLSPServer(ctx context.Context, rwc *ReadWriteCloser) {
	stream := jsonrpc2.NewBufferedStream(rwc, jsonrpc2.VSCodeObjectCodec{})
	c.conn = jsonrpc2.NewConn(ctx, stream, jsonrpc2.HandlerWithError(c.handleServerRequest))
	c.Logger.Infow("Connected to nxls server")
}

// handleServerRequest handles incoming requests from the server.
func (c *Client) handleServerRequest(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	if req.Notif && req.Method == "window/logMessage" {
		if c.isVerbose {
			params := &windowLogMessageNotification{}
			err := json.Unmarshal(*req.Params, params)
			if err != nil {
				return nil, err
			}

			c.Logger.Info(params.Message)
		}

		if c.isVerbose {
			c.Logger.Info(req)
		}

		return nil, nil
	}

	// Handle incoming requests from the server
	// You can implement your logic here
	return nil, nil
}
