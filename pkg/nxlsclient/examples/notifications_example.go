// Package examples demonstrates usage patterns for the nxlsclient package.
package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lazyengs/pkg/nxlsclient"
	"github.com/lazyengs/pkg/nxlsclient/commands"
	"go.lsp.dev/protocol"
)

// This example demonstrates how to use the notification listener
// to register and handle notifications from the Nx Language Server.
func ExampleNotificationHandling() {
	// Create a new client
	client := nxlsclient.NewClient("/path/to/nx/workspace", true)

	// 1. Basic notification handling
	// Register a handler for refresh workspace started notification
	refreshDisposable := client.OnNotification(
		commands.RefreshWorkspaceNotificationMethod,
		func(method string, params json.RawMessage) error {
			fmt.Println("Refresh workspace started!")
			return nil
		},
	)

	// 2. Typed notification handling with generic helper
	// Register a handler for window/logMessage with type checking
	logDisposable := client.OnNotification(
		nxlsclient.WindowLogMessageMethod,
		nxlsclient.TypedNotificationHandler(
			func(method string, params *nxlsclient.WindowLogMessage) error {
				// Now we have a properly typed parameter
				if params.Type == nxlsclient.LogError {
					fmt.Printf("Error from server: %s\n", params.Message)
				} else {
					fmt.Printf("Log from server: %s\n", params.Message)
				}
				return nil
			},
		),
	)

	// Start the client
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Initialize channel for receiving the initialization result
	initCh := make(chan *commands.InitializeRequestResult)

	// Start client in a goroutine
	go func() {
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
		err := client.Start(ctx, params, initCh)
		if err != nil {
			fmt.Printf("Client error: %v\n", err)
		}
	}()

	// Wait for initialization
	<-initCh
	fmt.Println("Client initialized successfully")

	// Use the client...

	// When you're done with specific notifications, dispose of handlers to avoid memory leaks
	refreshDisposable.Dispose()
	logDisposable.Dispose()

	// Stop the client when done (this also clears all remaining handlers)
	client.Stop(ctx)
}

// Example showing how to handle multiple notification types
func ExampleMultipleNotifications() {
	client := nxlsclient.NewClient("/path/to/nx/workspace", true)

	// Register a done channel to be notified when workspace refresh is complete
	// This pattern is useful for synchronizing with asynchronous operations
	refreshDone := make(chan struct{})

	// Track if we've seen the "started" notification
	var refreshStarted bool

	// Handler for refresh started
	client.OnNotification(
		commands.RefreshWorkspaceStartedNotificationMethod,
		func(method string, params json.RawMessage) error {
			fmt.Println("Refresh started...")
			refreshStarted = true
			return nil
		},
	)

	// Handler for refresh finished (regular notification)
	client.OnNotification(
		commands.RefreshWorkspaceNotificationMethod,
		func(method string, params json.RawMessage) error {
			fmt.Println("Refresh completed!")
			// Only signal completion if we saw the start notification
			if refreshStarted {
				close(refreshDone)
			}
			return nil
		},
	)

	// Start client...
	// trigger refresh...

	// Wait for refresh to complete with timeout
	select {
	case <-refreshDone:
		fmt.Println("Workspace refresh completed successfully")
	case <-time.After(30 * time.Second):
		fmt.Println("Timed out waiting for workspace refresh")
	}

	// Use client...

	// Cleanup
	client.Stop(context.Background())
}
