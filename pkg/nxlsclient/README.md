# nxlsclient

A Go client library for the Nx Language Server Protocol (LSP) server.

[![Go Reference](https://pkg.go.dev/badge/github.com/lazyengs/lazynx/pkg/nxlsclient.svg)](https://pkg.go.dev/github.com/lazyengs/lazynx/pkg/nxlsclient)

## Overview

`nxlsclient` provides a Go-native interface to the Nx Console Language Server Protocol (LSP) server. It enables Go applications to leverage the power of [Nx Console](https://github.com/nrwl/nx-console)'s tooling by:

- Managing the nxls server lifecycle (installation, startup, shutdown)
- Establishing a JSON-RPC connection to communicate with the server
- Providing a comprehensive API for interacting with Nx workspaces
- Supporting event notifications from the LSP server

This library is designed to work seamlessly with Nx workspace monorepos and integrates the nxls server directly, removing the need for external dependencies.

## Features

- **Zero external dependencies**: The nxls server is embedded directly in the library
- **Rich command API**: Full support for all Nx LSP commands
- **Event notifications**: Support for LSP notifications with typed handlers
- **Context-aware**: Full support for Go context for lifecycle management
- **Proper logging**: Structured logging with configurable verbosity
- **Cross-platform**: Works on macOS, Linux, and Windows

## Installation

```bash
go get github.com/lazyengs/lazynx/pkg/nxlsclient
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/lazyengs/lazynx/pkg/nxlsclient"
	"github.com/lazyengs/lazynx/pkg/nxlsclient/commands"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

func main() {
	// Setup logger
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()
	defer sugar.Sync()

	// Get absolute path to your Nx workspace
	nxWorkspacePath, err := filepath.Abs(".")
	if err != nil {
		sugar.Fatal(err)
	}

	// Create a new client
	client := nxlsclient.NewClient(nxWorkspacePath, true)
	
	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Handle termination signals
	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Create channel for initialization results
	ch := make(chan *commands.InitializeRequestResult)
	
	// Start the client asynchronously
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

	// Handle termination
	go func() {
		<-signalChan
		sugar.Info("Received interrupt signal")
		client.Stop(ctx)
		cancel()
		signal.Stop(signalChan)
	}()

	// Process initialization result
	init, ok := <-ch
	if ok {
		sugar.Infow("LSP server initialized successfully", "capabilities", init.Capabilities)
		
		// Now you can use the client.Commander to send requests
		workspace, err := client.Commander.SendWorkspaceRequest(ctx, &commands.WorkspaceRequestParams{
			Reset: false,
		})
		if err != nil {
			sugar.Errorf("Failed to get workspace: %v", err)
		} else {
			sugar.Infow("Retrieved workspace information", "version", workspace.NxVersion)
		}
	}

	// Wait for context to be done
	<-ctx.Done()
}
```

## Advanced Usage

### Command API

The client provides a comprehensive set of commands to interact with the Nx workspace through the `Commander` interface:

```go
// After client initialization
projectGraph, err := client.Commander.SendCreateProjectGraphRequest(ctx, commands.CreateProjectGraphParams{
    ShowAffected: false,
})
if err != nil {
    log.Fatalf("Failed to get project graph: %v", err)
}

fmt.Printf("Project graph contains %d nodes\n", len(projectGraph.Nodes))
```

### Notifications

The client supports event-based programming through LSP notifications:

```go
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
logDisposable.Dispose()
```

### Available Commands

The client supports all Nx LSP commands including:

- Workspace information
- Project graph generation
- Generator execution
- Target parsing and execution
- Project configuration
- Cloud integration

For a complete list of available commands, refer to the `commands` package documentation.

## Type Definitions

The library includes Nx-specific type definitions in the `nx-types` package:

- Project configuration
- Project graph models
- Generator schemas
- Workspace configuration
- Cloud integration types

## Examples

For more examples, check the `examples` directory in the package:

- [Notification Handling](./examples/notifications_example.go)

## Related Projects

- [Nx Console](https://github.com/nrwl/nx-console): The original Nx Console extension for VS Code
- [lazynx](https://github.com/lazyengs/lazynx): TUI application for Nx workspaces using this client

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## License

[MIT](LICENSE)