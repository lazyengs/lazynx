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

	"github.com/lazyengs/pkg/nxlsclient"
	"github.com/lazyengs/pkg/nxlsclient/commands"
	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
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
		fmt.Printf("Nx version: %s\n", workspace.NxVersion)
		fmt.Printf("Package manager: %s\n", workspace.PackageManager)
	}

	// 2. Get project graph
	projectGraph, err := client.Commander.SendCreateProjectGraphRequest(ctx, commands.CreateProjectGraphParams{
		ShowAffected: false,
	})
	if err != nil {
		fmt.Printf("Failed to get project graph: %v\n", err)
	} else {
		fmt.Printf("Project graph contains %d projects\n", len(projectGraph.Nodes))
		for name := range projectGraph.Nodes {
			fmt.Printf("- Project: %s\n", name)
		}
	}

	// 3. Get project by path
	currentDir, _ := os.Getwd()
	project, err := client.Commander.SendProjectByPathRequest(ctx, commands.ProjectByPathParams{
		ProjectPath: currentDir,
	})
	if err != nil {
		fmt.Printf("Failed to get project by path: %v\n", err)
	} else if project != nil {
		fmt.Printf("Current directory is part of project: %s\n", project.Name)
		fmt.Printf("Project root: %s\n", project.Root)
		fmt.Printf("Project type: %s\n", project.ProjectType)
	} else {
		fmt.Println("Current directory is not part of any project")
	}

	// 4. Get available generators
	generators, err := client.Commander.SendGeneratorsRequest(ctx, commands.GeneratorsParams{})
	if err != nil {
		fmt.Printf("Failed to get generators: %v\n", err)
	} else {
		fmt.Printf("Available generators: %d\n", len(generators))
		for i, gen := range generators {
			if i < 5 { // Only print first 5 generators
				fmt.Printf("- %s: %s\n", gen.Name, gen.Description)
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
func ExampleNotificationHandling() {
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
	projectGraph, err := client.Commander.SendCreateProjectGraphRequest(ctx, commands.CreateProjectGraphParams{
		ShowAffected: false,
	})
	if err != nil {
		fmt.Printf("Failed to get project graph: %v\n", err)
		return
	}

	// Print basic project graph information
	fmt.Printf("Project graph contains %d projects and %d external nodes\n",
		len(projectGraph.Nodes), len(projectGraph.ExternalNodes))

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

	// Sort projects by number of dependencies (implementation omitted for brevity)
	// ...

	// Print projects with the most dependencies
	fmt.Println("Projects with most dependencies:")
	for i, p := range projectsWithDeps {
		if i < 5 { // Only print top 5
			fmt.Printf("- %s: %d dependencies\n", p.name, p.deps)
		} else {
			break
		}
	}

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