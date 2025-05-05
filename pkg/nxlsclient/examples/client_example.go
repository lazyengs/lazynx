package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/lazyengs/lazynx/pkg/nxlsclient"
	"github.com/lazyengs/lazynx/pkg/nxlsclient/commands"
	nxtypes "github.com/lazyengs/lazynx/pkg/nxlsclient/nx-types"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

// This example demonstrates basic usage of the nxlsclient package.
// It shows how to create a client, initialize it, and perform common operations.
func ExampleBasicUsage() {
	// Setup logger
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()
	defer sugar.Sync()

	// Get absolute path to your Nx workspace (use the current directory in this example)
	nxWorkspacePath, err := filepath.Abs(".")
	if err != nil {
		sugar.Fatal(err)
	}

	// Create a new client
	client := nxlsclient.NewClient(nxWorkspacePath, true)

	// Setup context with cancellation and timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Setup signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		fmt.Println("Received interrupt signal, shutting down...")
		cancel()
		signal.Stop(signalChan)
	}()

	// Create channel for initialization results
	initCh := make(chan *commands.InitializeRequestResult)

	// Start client in a goroutine
	go func() {
		// Setup initialization parameters
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

		// Start the client
		err := client.Start(ctx, params, initCh)
		if err != nil {
			fmt.Printf("Client error: %v\n", err)
		}
	}()

	// Wait for initialization
	init, ok := <-initCh
	if !ok {
		fmt.Println("Failed to initialize client")
		return
	}

	fmt.Println("Client initialized successfully!")
	fmt.Printf("Server capabilities: %+v\n", init.Capabilities)

	// Use the client to interact with the Nx workspace
	// 1. Get workspace information
	workspace, err := client.Commander.SendWorkspaceRequest(ctx, &commands.WorkspaceRequestParams{
		Reset: false,
	})
	if err != nil {
		fmt.Printf("Failed to get workspace information: %v\n", err)
	} else {
		fmt.Printf("Nx version: %s\n", workspace.NxVersion.Full)
		// Package manager is not directly available in NxWorkspace
		fmt.Printf("Is Lerna: %t\n", workspace.IsLerna)
	}

	// 2. Get project graph
	projectGraphStr, err := client.Commander.SendCreateProjectGraphRequest(ctx, commands.CreateProjectGraphParams{
		ShowAffected: false,
	})
	if err != nil {
		fmt.Printf("Failed to get project graph: %v\n", err)
	} else if projectGraphStr != nil {
		// Project graph is returned as a string that would need to be parsed
		// In a real implementation, you would parse the JSON string into a ProjectGraph object
		fmt.Printf("Got project graph (length: %d chars)\n", len(*projectGraphStr))
		// For demonstration purposes only:
		fmt.Printf("Project graph snippet: %.100s...\n", *projectGraphStr)
	}

	// 3. Get project by path
	currentDir, _ := os.Getwd()
	project, err := client.Commander.SendProjectByPathRequest(ctx, commands.ProjectByPathParams{
		ProjectPath: currentDir,
	})
	if err != nil {
		fmt.Printf("Failed to get project by path: %v\n", err)
	} else if project != nil {
		// Need to check pointer fields before accessing
		projectName := "unnamed"
		if project.Name != nil {
			projectName = *project.Name
		}

		projectTypeStr := "unknown"
		if project.ProjectType != nil {
			projectTypeStr = string(*project.ProjectType)
		}

		fmt.Printf("Current directory is part of project: %s\n", projectName)
		fmt.Printf("Project root: %s\n", project.Root)
		fmt.Printf("Project type: %s\n", projectTypeStr)
	} else {
		fmt.Println("Current directory is not part of any project")
	}

	// 4. Get available generators
	generators, err := client.Commander.SendGeneratorsRequest(ctx, commands.GeneratorsRequestParams{})
	if err != nil {
		fmt.Printf("Failed to get generators: %v\n", err)
	} else {
		fmt.Printf("Available generators: %d\n", len(generators))
		for i, gen := range generators {
			if i < 5 { // Only print first 5 generators
				// Access the description from the data field if available
				description := ""
				if gen.Data != nil {
					description = gen.Data.Description
				}
				fmt.Printf("- %s: %s\n", gen.Name, description)
			}
		}
		if len(generators) > 5 {
			fmt.Printf("  ... and %d more\n", len(generators)-5)
		}
	}

	// Clean up and stop the client
	client.Stop(ctx)
	fmt.Println("Client stopped successfully")
}

// This example demonstrates working with LSP notifications from the Nx server.
// Example renamed to avoid duplicate with notifications_example.go
func ExampleHandlingClientNotifications() {
	// Create a new client
	client := nxlsclient.NewClient("/path/to/nx/workspace", true)

	// Create notification channels to communicate with the main goroutine
	refreshStarted := make(chan struct{}, 1)
	refreshCompleted := make(chan struct{}, 1)
	logMessages := make(chan string, 10)

	// Register handlers for different notification types

	// 1. Register for refresh workspace started notification
	client.OnNotification(
		commands.RefreshWorkspaceStartedNotificationMethod,
		func(method string, params json.RawMessage) error {
			fmt.Println("Refresh workspace started!")
			refreshStarted <- struct{}{}
			return nil
		},
	)

	// 2. Register for refresh workspace completed notification
	client.OnNotification(
		nxlsclient.NxRefreshWorkspaceMethod,
		func(method string, params json.RawMessage) error {
			fmt.Println("Refresh workspace completed!")
			refreshCompleted <- struct{}{}
			return nil
		},
	)

	// 3. Register for log messages with type checking
	client.OnNotification(
		nxlsclient.WindowLogMessageMethod,
		nxlsclient.TypedNotificationHandler(
			func(method string, params *nxlsclient.WindowLogMessage) error {
				var level string
				switch params.Type {
				case nxlsclient.LogError:
					level = "ERROR"
				case nxlsclient.LogWarning:
					level = "WARNING"
				case nxlsclient.LogInfo:
					level = "INFO"
				case nxlsclient.LogDebug:
					level = "DEBUG"
				}

				message := fmt.Sprintf("[%s] %s", level, params.Message)
				fmt.Println(message)
				logMessages <- message
				return nil
			},
		),
	)

	// Start the client (implementation omitted for brevity)
	// ...

	// Example of waiting for a workspace refresh to complete
	func() {
		// Request a workspace refresh
		ctx := context.Background()
		if err := client.Commander.SendWorkspaceRefreshNotification(ctx); err != nil {
			fmt.Printf("Failed to send refresh notification: %v\n", err)
			return
		}

		// Wait for both start and completion notifications
		select {
		case <-refreshStarted:
			fmt.Println("Received refresh started notification")
		case <-time.After(5 * time.Second):
			fmt.Println("Timed out waiting for refresh to start")
			return
		}

		select {
		case <-refreshCompleted:
			fmt.Println("Received refresh completed notification")
		case <-time.After(30 * time.Second):
			fmt.Println("Timed out waiting for refresh to complete")
			return
		}

		fmt.Println("Workspace refresh cycle completed successfully")
	}()

	// Example of collecting log messages for a period of time
	func() {
		// Collect logs for 5 seconds
		timeout := time.After(5 * time.Second)
		var logs []string

		// Collect logs until timeout
	collectLoop:
		for {
			select {
			case log := <-logMessages:
				logs = append(logs, log)
			case <-timeout:
				break collectLoop
			}
		}

		fmt.Printf("Collected %d log messages\n", len(logs))
		for i, log := range logs {
			if i < 5 { // Only print first 5 logs
				fmt.Printf("%d: %s\n", i+1, log)
			}
		}
		if len(logs) > 5 {
			fmt.Printf("  ... and %d more\n", len(logs)-5)
		}
	}()

	// Clean up
	client.Stop(context.Background())
}

// This example demonstrates how to get and process the project graph.
func ExampleProjectGraph() {
	// Create and initialize client
	client := nxlsclient.NewClient("/path/to/nx/workspace", true)
	ctx := context.Background()

	// Start client and wait for initialization (implementation omitted for brevity)
	// ...

	// Get project graph
	projectGraphStr, err := client.Commander.SendCreateProjectGraphRequest(ctx, commands.CreateProjectGraphParams{
		ShowAffected: false,
	})
	if err != nil {
		fmt.Printf("Failed to get project graph: %v\n", err)
		return
	}
	if projectGraphStr == nil {
		fmt.Println("No project graph returned")
		return
	}

	// In a real implementation, you would parse the JSON string into a ProjectGraph object
	// For example:
	// var projectGraph nxtypes.ProjectGraph
	// if err := json.Unmarshal([]byte(*projectGraphStr), &projectGraph); err != nil {
	//     fmt.Printf("Failed to parse project graph: %v\n", err)
	//     return
	// }

	// For this example, we'll just demonstrate how you might work with a ProjectGraph
	// after parsing it properly
	fmt.Printf("Got project graph data (length: %d chars)\n", len(*projectGraphStr))
	fmt.Println("In a real implementation, you would:")
	fmt.Println("1. Parse the JSON string into a ProjectGraph object")
	fmt.Println("2. Access graph properties like Nodes, ExternalNodes, and Dependencies")
	fmt.Println("3. Process the graph data for your specific use case")

	// This is just example code that would work with a properly parsed graph:
	/*
		// Print basic project graph information
		nodeCount := len(projectGraph.Nodes)
		externalNodeCount := len(projectGraph.ExternalNodes)
		fmt.Printf("Project graph contains %d projects and %d external nodes\n",
			nodeCount, externalNodeCount)

		// Count projects by type
		projectTypeCount := make(map[string]int)
		for _, node := range projectGraph.Nodes {
			projectTypeCount[node.Type]++
		}
		fmt.Println("Projects by type:")
		for typ, count := range projectTypeCount {
			fmt.Printf("- %s: %d\n", typ, count)
		}

		// Find projects with the most dependencies
		type projectWithDeps struct {
			name string
			deps int
		}
		var projectsWithDeps []projectWithDeps
		for name, deps := range projectGraph.Dependencies {
			projectsWithDeps = append(projectsWithDeps, projectWithDeps{name, len(deps)})
		}
	*/

	// In a complete implementation, you would:
	// 1. Sort projects by number of dependencies
	// 2. Print the projects with the most dependencies

	fmt.Println("Projects with most dependencies would be listed here")
	fmt.Println("For example:")
	fmt.Println("- project1: 15 dependencies")
	fmt.Println("- project2: 12 dependencies")
	fmt.Println("- project3: 8 dependencies")

	// Find projects affected by a specific file
	affectedProjects := make(map[string]*nxtypes.ProjectConfiguration)
	targetFile := "/path/to/some/file.ts"

	// This is a simplified approach - in a real implementation, you would use:
	// 1. SendSourceMapFilesToProjectsMapRequest to get the file-to-project mapping
	// 2. Or derive affected projects from the dependency graph

	fmt.Printf("Projects affected by changes to %s:\n", targetFile)
	for name, project := range affectedProjects {
		fmt.Printf("- %s (%s)\n", name, project.Root)
	}

	// Clean up
	client.Stop(ctx)
}
