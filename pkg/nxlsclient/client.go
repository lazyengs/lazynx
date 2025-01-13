package nxlsclient

import (
	"context"
	"time"

	"github.com/sourcegraph/jsonrpc2"
	"go.uber.org/zap"
)

type Client struct {
	Logger          *zap.SugaredLogger
	conn            *jsonrpc2.Conn
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
func (c *Client) Start(ctx context.Context) error {
	c.Logger.Infow("Starting client")

	err := c.unpackServer()
	if err != nil {
		c.stop()
		return err
	}
	err = c.installDependencies(ctx)
	if err != nil {
		c.stop()
		return err
	}

	rwc, err := c.startNxls(ctx)
	if err != nil {
		c.stop()
		return err
	}

	c.connectToLSPServer(ctx, rwc)

	err = c.sendInitializeCommand(ctx)
	if err != nil {
		c.stop()
		return err
	}

	<-ctx.Done()
	c.stop()

	return nil
}

// stop gracefully stops the client, cleaning up resources and closing connections.
func (c *Client) stop() {
	c.Logger.Infow("Stopping client")
	c.Logger.Sync()
	c.cleanUpServer()
	c.conn.Close()
	// grace time to allow the completion of the clean up
	time.Sleep(1 * time.Second)
}
