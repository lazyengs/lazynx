package nxlsclient

import (
	"context"

	"github.com/lazyengs/pkg/nxlsclient/commands"
	"github.com/sourcegraph/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

type Client struct {
	Logger               *zap.SugaredLogger
	conn                 *jsonrpc2.Conn
	serverDir            string
	NxWorkspacePath      string
	isVerbose            bool
	Commander            *commands.Commander
	notificationListener *notificationListener
}

// NewClient creates a new Client struct instance with the given nxWorkspacePath and verbosity level.
func NewClient(nxWorkspacePath string, verbose bool) *Client {
	logger, _ := zap.NewDevelopment()
	if !verbose {
		logger, _ = zap.NewProduction()
	}
	sugar := logger.Sugar()

	sugar.Debugw("Creating new client")

	return &Client{
		Logger:               sugar,
		NxWorkspacePath:      nxWorkspacePath,
		isVerbose:            verbose,
		notificationListener: newNotificationListener(),
	}
}

// Start spawns the nxls server process, sends the initialize command to the LSP server and listen for incoming messages.
func (c *Client) Start(ctx context.Context, initParams *protocol.InitializeParams, ch chan *commands.InitializeRequestResult) error {
	c.Logger.Debugw("Starting client")

	err := c.unpackServer()
	if err != nil {
		c.Stop(ctx)
		return err
	}
	err = c.installDependencies(ctx)
	if err != nil {
		c.Stop(ctx)
		return err
	}

	rwc, err := c.startNxls(ctx)
	if err != nil {
		c.Stop(ctx)
		return err
	}

	c.connectToLSPServer(ctx, rwc)

	c.Commander = commands.NewCommander(c.conn, c.Logger)

	initResponse, err := c.Commander.SendInitializeRequest(ctx, initParams)

	ch <- initResponse
	close(ch)
	if err != nil {
		c.Stop(ctx)
		return err
	}

	<-ctx.Done()

	return nil
}

// Stop gracefully Stops the client, cleaning up resources and closing connections.
func (c *Client) Stop(ctx context.Context) {
	c.Logger.Debugw("Stopping client")

	// Clear all notification handlers
	if c.notificationListener != nil {
		c.notificationListener.clearHandlers()
	}

	err := c.stopNxls(ctx)
	if err != nil {
		c.Logger.Errorw("An error occurred while stopping nxls", "error", err.Error())
	}

	c.Logger.Debugw("Clean up completed")
	c.Logger.Sync()
}

// OnNotification registers a handler for a specific notification method.
// Returns a Disposable that can be used to unregister the handler.
func (c *Client) OnNotification(method string, handler NotificationHandler) *Disposable {
	if c.notificationListener == nil {
		c.Logger.Warnw("Notification listener is nil, creating a new one")
		c.notificationListener = newNotificationListener()
	}
	return c.notificationListener.registerHandler(method, handler)
}
