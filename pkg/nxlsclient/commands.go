package nxlsclient

import (
	"context"
)

func (c *Client) sendRequest(ctx context.Context, method string, params ...any) (any, error) {
	var result any
	c.Logger.Debugw("Sending request", "method", method, "params", params)

	if err := c.conn.Call(ctx, method, params, &result); err != nil {
		c.Logger.Errorw("An error occurred while executing the request",
			"method", method, "params", params,
			"error", err.Error(),
		)
		return nil, err
	}

	return result, nil
}

func (c *Client) sendNotification(ctx context.Context, method string, params []any) error {
	c.Logger.Debugw("Sending notification", method, "params", params)

	if err := c.conn.Notify(ctx, method, params); err != nil {
		c.Logger.Errorw("An error occurred while sending the notification",
			"method", method, "params", params,
			"error", err.Error(),
		)
		return err
	}

	return nil
}
