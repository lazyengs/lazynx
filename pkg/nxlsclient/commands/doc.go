/*
Package commands provides implementations of all Nx Language Server Protocol (LSP) commands.

This package contains request and notification implementations for interacting with an Nx workspace
through the Nx LSP server. Each command is implemented as a method on the Commander type.

# Overview

The commands package provides:

  - Request implementations for retrieving workspace information
  - Project graph generation and manipulation
  - Generator execution and configuration
  - Target parsing and execution
  - File-to-project mappings
  - Cloud integration commands
  - Notification handlers for workspace events

# Usage

Commands are typically used through the Commander instance provided by the nxlsclient package:

	// After client initialization
	workspace, err := client.Commander.SendWorkspaceRequest(ctx, &commands.WorkspaceRequestParams{
		Reset: false,
	})

	if err != nil {
		// Handle error
	}

	fmt.Printf("Nx version: %s\n", workspace.NxVersion)

# Available Commands

The package includes commands for:

  - Workspace: Retrieving and serializing workspace configuration
  - Project Graph: Creating and navigating project dependency graphs
  - Projects: Finding projects by path, root, or other criteria
  - Generators: Executing generators with various options
  - Targets: Parsing and executing targets within projects
  - Cloud: Interacting with Nx Cloud services
  - File Operations: Mapping files to projects

# Custom Commander

You can create a custom Commander instance with the NewCommander function:

	conn := // your jsonrpc2.Conn instance
	logger := // your zap.SugaredLogger instance
	commander := commands.NewCommander(conn, logger)

	// Now use the commander to send requests
	result, err := commander.SendWorkspaceRequest(ctx, params)

# Error Handling

All command methods return errors that should be handled by the caller:

	result, err := client.Commander.SendCreateProjectGraphRequest(ctx, params)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			// Handle cancellation
		case strings.Contains(err.Error(), "connection is nil"):
			// Handle connection problems
		default:
			// Handle other errors
		}
	}
*/
package commands
