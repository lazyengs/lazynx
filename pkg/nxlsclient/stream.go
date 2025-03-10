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
	c.Logger.Debugw("Connected to nxls server")
}

// handleServerRequest handles incoming requests from the server.
func (c *Client) handleServerRequest(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	// If notification, pass to registered handlers
	if req.Notif {
		// Special case for window/logMessage
		if req.Method == "window/logMessage" {
			if c.isVerbose {
				params := &windowLogMessageNotification{}
				err := json.Unmarshal(*req.Params, params)
				if err != nil {
					return nil, err
				}

				c.Logger.Info(params.Message)
			}
		}

		// Log all notifications when verbose is enabled
		if c.isVerbose {
			c.Logger.Infow("Received notification", "method", req.Method)
		}

		// Check if we have handlers for this notification method
		if c.notificationListener != nil && c.notificationListener.hasHandlers(req.Method) {
			// Process asynchronously to avoid blocking the JSONRPC handler
			go func(method string, params json.RawMessage) {
				c.Logger.Debugw("Notifying handlers", "method", method)
				c.notificationListener.notifyAll(method, params)
			}(req.Method, *req.Params)
		}

		return nil, nil
	}

	// Handle non-notification requests (actual RPC calls)
	if c.isVerbose {
		c.Logger.Infow("Received request", "method", req.Method, "id", req.ID)
	}

	// For now, we don't handle any requests from the server
	return nil, nil
}
