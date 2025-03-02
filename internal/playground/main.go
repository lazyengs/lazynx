package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/lazyengs/pkg/nxlsclient"
	"github.com/lazyengs/pkg/nxlsclient/commands"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

func main() {
	_logger, _ := zap.NewDevelopment()
	logger := _logger.Sugar()
	defer logger.Sync()

	logger.Infow("Starting Nx client playground")

	// Get the workspace path
	currentNxWorkspacePath, err := filepath.Abs("../..")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infow("Using workspace path", "path", currentNxWorkspacePath)

	// Create the client with verbose logging
	client := nxlsclient.NewClient(currentNxWorkspacePath, true)

	// Create a cancellable context for the main app
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure context is cancelled when main exits

	// Ctrl+c like signal detection
	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Channel for initialization result
	ch := make(chan *commands.InitializeRequestResult)

	// Start client in a goroutine
	go func() {
		logger.Infow("Starting client...")
		params := &protocol.InitializeParams{
			RootURI: protocol.DocumentURI(client.NxWorkspacePath),
			Capabilities: protocol.ClientCapabilities{
				Workspace: &protocol.WorkspaceClientCapabilities{
					Configuration: true,
				},
				TextDocument: &protocol.TextDocumentClientCapabilities{},
			},
			InitializationOptions: map[string]any{
				"workspacePath": client.NxWorkspacePath,
			},
		}
		err := client.Start(ctx, params, ch)
		if err != nil {
			logger.Errorw("Error starting client", "error", err)
		}
	}()

	// Setup signal handler for graceful shutdown
	go func() {
		sig := <-signalChan
		logger.Infow("Received signal", "signal", sig)

		logger.Infow("Starting shutdown sequence")
		// Stop the client - this will also stop the daemon
		client.Stop(ctx)
		logger.Infow("Shutdown sequence completed")

		// Cancel the main context to stop everything
		cancel()
		signal.Stop(signalChan)

		// Ensure clean exit
		logger.Sync()
		os.Exit(0)
	}()

	// Wait for initialization to complete
	logger.Infow("Waiting for initialization...")
	init, ok := <-ch
	if !ok {
		logger.Errorw("Initialization channel closed unexpectedly")
		return
	}
	logger.Infow("Initialization complete", "capabilities", init.Capabilities)

	// Get the commander for sending requests
	logger.Infow("Executing commands...")
	commander := client.Commander
	if commander == nil {
		logger.Errorw("Commander is nil, cannot execute commands")
		return
	}

	// Example request
	result, err := commander.SendWorkspaceRequest(ctx, &commands.WorkspaceRequestParams{
		Reset: false,
	})
	if err != nil {
		logger.Errorw("Error sending workspace request", "error", err)
	} else {
		logger.Infow("Workspace command successful")

		// Output some interesting details from the result
		if result != nil {
			if proj, ok := result.ProjectGraph.Nodes["pkg/nxlsclient"]; ok {
				logger.Infow("Package nxlsclient details", "data", proj)
			}
		}
	}

	logger.Infow("Playground running, press Ctrl+C to stop")
	<-ctx.Done()
	logger.Infow("Main context done, exiting")
}
