package nxlsclient

import (
	"context"

	"github.com/sourcegraph/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

type Client struct {
	Logger          *zap.SugaredLogger
	conn            *jsonrpc2.Conn
	serverDir       string
	NxWorkspacePath string
	isVerbose       bool
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
		Logger:          sugar,
		NxWorkspacePath: nxWorkspacePath,
		isVerbose:       verbose,
	}
}

// Start spawns the nxls server process, sends the initialize command to the LSP server and listen for incoming messages.
func (c *Client) Start(ctx context.Context, initParams *protocol.InitializeParams, ch chan *InitializeCommandResult) error {
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

	initResponse, err := c.sendInitializeCommand(ctx, initParams)
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

	c.stopNxls(ctx)

	if err := c.cleanUpServer(); err != nil {
		c.Logger.Errorw("An error occurred while cleaning up the server", "error", err.Error())
	}
	c.Logger.Debugw("Clean up completed")
	c.Logger.Sync()
}
