/*
Package nxlsclient provides a Go client for interacting with the Nx Language Server Protocol (LSP) server.

The nxlsclient package enables Go applications to leverage Nx Console's powerful tooling by embedding
the Nx language server directly in the package and providing a clean API for communication.

# Overview

The nxlsclient's primary features include:

  - Managing the LSP server lifecycle (installation, startup, shutdown)
  - Maintaining a JSON-RPC connection for command execution
  - Handling LSP notifications with a flexible event system
  - Providing typed access to Nx workspace data structures
  - Supporting all Nx commands through a comprehensive Commander API

# Basic Usage

Creating a client and initializing the LSP connection:

	client := nxlsclient.NewClient("/path/to/nx/workspace", true)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize channel for receiving the initialization result
	ch := make(chan *commands.InitializeRequestResult)

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
		client.Start(ctx, params, ch)
	}()

	// Wait for initialization
	init := <-ch

	// When done
	client.Stop(ctx)

# Commands

After initialization, you can use the Commander to interact with the Nx workspace:

	// Get workspace information
	workspace, err := client.Commander.SendWorkspaceRequest(ctx, &commands.WorkspaceRequestParams{
		Reset: false,
	})

	// Get project graph
	projectGraph, err := client.Commander.SendCreateProjectGraphRequest(ctx, commands.CreateProjectGraphParams{
		ShowAffected: false,
	})

# Notifications

The client supports event-based programming through LSP notifications:

	// Register for workspace refresh notifications
	refreshDisposable := client.OnNotification(
		nxlsclient.NxRefreshWorkspaceMethod,
		func(method string, params json.RawMessage) error {
			fmt.Println("Workspace refresh completed!")
			return nil
		},
	)

	// Register for log message notifications with type checking
	logDisposable := client.OnNotification(
		nxlsclient.WindowLogMessageMethod,
		nxlsclient.TypedNotificationHandler(
			func(method string, params *nxlsclient.WindowLogMessage) error {
				if params.Type == nxlsclient.LogError {
					fmt.Printf("Error from server: %s\n", params.Message)
				}
				return nil
			},
		),
	)

	// Clean up when done
	refreshDisposable.Dispose()

# Related Packages

  - commands: Contains all LSP command implementations
  - nx-types: Contains Nx-specific type definitions
*/
package nxlsclient

