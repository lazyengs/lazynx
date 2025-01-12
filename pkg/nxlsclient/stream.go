package nxlsclient

// func (c *Client) listenToLSPMessages(ctx context.Context, rwc *ReadWriteCloser) {
//
// 	c.conn.Go(ctx, func(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
// 		c.logger.Debugw("Received message", "method", req.Method(), "params", req.Params())
//
// 		if req.Method() == protocol.MethodWindowLogMessage {
// 			var params protocol.LogMessageParams
// 			if err := json.Unmarshal(req.Params(), &params); err == nil {
// 				c.logger.Info(params.Message)
// 			}
// 		} else {
// 			select {
// 			case c.notifications <- Notification{Method: req.Method(), Params: req.Params()}:
// 			default:
// 				c.logger.Warn("Notification channel full, dropping message")
// 			}
// 		}
// 		return nil
// 	})
// }
