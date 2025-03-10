package commands

import (
	"context"
	"fmt"

	"github.com/sourcegraph/jsonrpc2"
	"go.uber.org/zap"
)

// Commander is responsible for sending requests and notifications via JSON-RPC.
type Commander struct {
	Logger *zap.SugaredLogger // Logger is used to log messages.
	conn   *jsonrpc2.Conn     // conn is the JSON-RPC connection.
}

// NewCommander creates a new Commander instance.
func NewCommander(conn *jsonrpc2.Conn, logger *zap.SugaredLogger) *Commander {
	return &Commander{
		Logger: logger,
		conn:   conn,
	}
}

// sendRequest sends a JSON-RPC request and stores the result in the provided result parameter.
func (c *Commander) sendRequest(ctx context.Context, method string, params any, result any) error {
	c.Logger.Debugw("Sending request", "method", method, "params", params)

	// Check connection state before making the call
	if c.conn == nil {
		return fmt.Errorf("connection is nil")
	}

	if err := c.conn.Call(ctx, method, params, result); err != nil {
		c.Logger.Warnw("Request failed", "method", method, "error", err)
		return fmt.Errorf("an error occurred while executing the request: %w", err)
	}

	c.Logger.Debugw("Request successful", "method", method)
	return nil
}

// sendNotification sends a JSON-RPC notification.
func (c *Commander) sendNotification(ctx context.Context, method string, params any) error {
	c.Logger.Debugw("Sending notification", "method", method, "params", params)

	// Check connection state before making the call
	if c.conn == nil {
		return fmt.Errorf("connection is nil")
	}

	if err := c.conn.Notify(ctx, method, params); err != nil {
		c.Logger.Warnw("Notification failed", "method", method, "error", err)
		return fmt.Errorf("an error occurred while sending the notification: %w", err)
	}

	c.Logger.Debugw("Notification sent successfully", "method", method)
	return nil
}
