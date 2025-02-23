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

	if err := c.conn.Call(ctx, method, params, &result); err != nil {
		c.Logger.Errorw("An error occurred while executing the request",
			"method", method, "params", params,
			"error", err.Error(),
		)
		return fmt.Errorf("an error occurred while executing the request: %w", err)
	}

	return nil
}

// sendNotification sends a JSON-RPC notification.
func (c *Commander) sendNotification(ctx context.Context, method string, params any) error {
	c.Logger.Debugw("Sending notification", "method", method, "params", params)

	if err := c.conn.Notify(ctx, method, params); err != nil {
		c.Logger.Errorw("An error occurred while sending the notification",
			"method", method, "params", params,
			"error", err.Error(),
		)
		return fmt.Errorf("an error occurred while sending the notification: %w", err)
	}

	return nil
}
