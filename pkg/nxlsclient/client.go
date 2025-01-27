package nxlsclient

import (
	"context"
	"os/exec"

	"github.com/sourcegraph/jsonrpc2"
	"go.uber.org/zap"
)

type Client struct {
	Logger          *zap.SugaredLogger
	conn            *jsonrpc2.Conn
	cmd             *exec.Cmd
	serverDir       string
	nxWorkspacePath string
	isVerbose       bool
}

// NewClient creates a new Client struct instance with the given nxWorkspacePath and verbosity level.
func NewClient(nxWorkspacePath string, verbose bool) *Client {
	logger, _ := zap.NewDevelopment()
	if !verbose {
		logger, _ = zap.NewProduction()
	}
	sugar := logger.Sugar()

	sugar.Infow("Creating new client")

	return &Client{
		Logger:          sugar,
		nxWorkspacePath: nxWorkspacePath,
		isVerbose:       verbose,
	}
}

// Start spawns the nxls server process, sends the initialize command to the LSP server and listen for incoming messages.
func (c *Client) Start(ctx context.Context, ch chan *InitializeCommandResult) error {
	c.Logger.Infow("Starting client")

	err := c.unpackServer()
	if err != nil {
		c.Stop()
		return err
	}
	err = c.installDependencies(ctx)
	if err != nil {
		c.Stop()
		return err
	}

	rwc, err := c.startNxls(ctx)
	if err != nil {
		c.Stop()
		return err
	}

	c.connectToLSPServer(ctx, rwc)

	initResponse, err := c.sendInitializeCommand(ctx)
	ch <- initResponse
	close(ch)
	if err != nil {
		c.Stop()
		return err

	}

	<-ctx.Done()

	return nil
}

// Stop gracefully Stops the client, cleaning up resources and closing connections.
func (c *Client) Stop() {
	c.Logger.Infow("Stopping client")
	if c.conn != nil {
		c.conn.Close()
	}
	c.Logger.Infow("Clean up completed")
	c.Logger.Sync()
}
