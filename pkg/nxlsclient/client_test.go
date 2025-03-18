package nxlsclient_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/lazyengs/lazynx/pkg/nxlsclient"
	"github.com/lazyengs/lazynx/pkg/nxlsclient/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.lsp.dev/protocol"
)

// TestClientE2E performs end-to-end tests of the nxlsclient's functionality
// by executing actual commands against a real NX workspace.
func TestClientE2E(t *testing.T) {
	// Skip the test if running in CI or if SKIP_E2E_TESTS is set
	if os.Getenv("CI") != "" || os.Getenv("SKIP_E2E_TESTS") != "" {
		t.Skip("Skipping E2E tests in CI environment or when SKIP_E2E_TESTS is set")
	}

	// Find the workspace path - use the project root directory
	workspacePath, err := filepath.Abs("../../..")
	require.NoError(t, err, "Failed to get absolute path")

	// Create a new client
	client := nxlsclient.NewClient(workspacePath, false)
	require.NotNil(t, client, "Client should not be nil")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize client
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
		err := client.Start(ctx, params, ch)
		if err != nil {
			t.Errorf("Failed to start client: %v", err)
		}
	}()

	// Wait for initialization
	var initResult *commands.InitializeRequestResult
	select {
	case initResult = <-ch:
		require.NotNil(t, initResult, "Initialization result should not be nil")
	case <-time.After(20 * time.Second):
		t.Fatal("Timeout waiting for initialization")
	}

	// Create a snapshotter instance with custom config
	updateSnapshots := os.Getenv("UPDATE_SNAPSHOTS") == "true"
	snapshotter := cupaloy.New(cupaloy.SnapshotSubdirectory("testdata/snapshots"))

	// Execute tests for different commands
	t.Run("InitializeResult", func(t *testing.T) {
		// Make a deep copy of the init result to normalize the PID
		initResultMap := make(map[string]interface{})

		// Convert to JSON and back to create a deep copy
		jsonData, err := json.Marshal(initResult)
		require.NoError(t, err, "Failed to marshal init result")

		err = json.Unmarshal(jsonData, &initResultMap)
		require.NoError(t, err, "Failed to unmarshal init result")

		// Normalize the PID to a fixed value to ensure snapshot consistency
		if pidValue, ok := initResultMap["pid"]; ok && pidValue != nil {
			initResultMap["pid"] = 12345 // Use a fixed pid for snapshot comparisons
		}

		// Snapshot the normalized initialization result
		prettyJson, err := prettyPrintJSON(initResultMap)
		require.NoError(t, err, "Failed to marshal init result to JSON")

		if updateSnapshots {
			snapshotter.SnapshotT(t, prettyJson)
			t.Log("Snapshot updated for InitializeResult")
		} else {
			// The SnapshotT function automatically fails the test if the snapshot doesn't match
			snapshotter.SnapshotT(t, prettyJson)
		}
	})

	// Skip if commander isn't available for all remaining tests
	if client.Commander == nil {
		t.Skip("Commander not available for remaining tests")
	}

	t.Run("WorkspaceRequest", func(t *testing.T) {
		// Execute workspace request
		result, err := client.Commander.SendWorkspaceRequest(ctx, &commands.WorkspaceRequestParams{
			Reset: false,
		})
		require.NoError(t, err, "Workspace request should not error")

		// Snapshot the workspace result
		prettyJson, err := prettyPrintJSON(result)
		require.NoError(t, err, "Failed to marshal workspace result to JSON")

		if updateSnapshots {
			snapshotter.SnapshotT(t, prettyJson)
			t.Log("Snapshot updated for WorkspaceRequest")
		} else {
			snapshotter.SnapshotT(t, prettyJson)
		}
	})

	t.Run("ProjectGraphRequest", func(t *testing.T) {
		// Execute project graph request
		result, err := client.Commander.SendCreateProjectGraphRequest(ctx, commands.CreateProjectGraphParams{
			ShowAffected: false,
		})
		if err != nil {
			// For this test, we allow errors since project graph might not be available in all environments
			t.Logf("Project graph request error (this may be expected): %v", err)
			snapshot := fmt.Sprintf("Error: %v", err)

			if updateSnapshots {
				snapshotter.SnapshotT(t, snapshot)
				t.Log("Snapshot updated for ProjectGraphRequest (error case)")
			} else {
				snapshotter.SnapshotT(t, snapshot)
			}
			return
		}

		// Snapshot the project graph result
		prettyJson, err := prettyPrintJSON(result)
		require.NoError(t, err, "Failed to marshal project graph result to JSON")

		if updateSnapshots {
			snapshotter.SnapshotT(t, prettyJson)
			t.Log("Snapshot updated for ProjectGraphRequest")
		} else {
			snapshotter.SnapshotT(t, prettyJson)
		}
	})

	t.Run("WorkspaceSerializedRequest", func(t *testing.T) {
		// Execute workspace serialized request
		result, err := client.Commander.SendWorkspaceSerializedRequest(ctx, commands.WorkspaceRequestParams{
			Reset: false,
		})
		require.NoError(t, err, "Workspace serialized request should not error")

		// Snapshot the workspace serialized result
		prettyJson, err := prettyPrintJSON(result)
		require.NoError(t, err, "Failed to marshal workspace serialized result to JSON")

		if updateSnapshots {
			snapshotter.SnapshotT(t, prettyJson)
			t.Log("Snapshot updated for WorkspaceSerializedRequest")
		} else {
			snapshotter.SnapshotT(t, prettyJson)
		}
	})

	t.Run("ProjectByPathRequest", func(t *testing.T) {
		// Execute project by path request
		result, err := client.Commander.SendProjectByPathRequest(ctx, commands.ProjectByPathParams{
			ProjectPath: filepath.Join(workspacePath, "pkg/nxlsclient"),
		})
		if err != nil {
			// For this test, we allow errors since project might not be available in all environments
			t.Logf("Project by path request error (this may be expected): %v", err)
			snapshot := fmt.Sprintf("Error: %v", err)

			if updateSnapshots {
				snapshotter.SnapshotT(t, snapshot)
				t.Log("Snapshot updated for ProjectByPathRequest (error case)")
			} else {
				snapshotter.SnapshotT(t, snapshot)
			}
			return
		}

		// Snapshot the project by path result
		prettyJson, err := prettyPrintJSON(result)
		require.NoError(t, err, "Failed to marshal project by path result to JSON")

		if updateSnapshots {
			snapshotter.SnapshotT(t, prettyJson)
			t.Log("Snapshot updated for ProjectByPathRequest")
		} else {
			snapshotter.SnapshotT(t, prettyJson)
		}
	})

	// Test the notifications (these don't have responses to snapshot)
	t.Run("RefreshWorkspaceNotification", func(t *testing.T) {
		// Execute refresh workspace notification
		err := client.Commander.SendWorkspaceRefreshNotification(ctx)
		assert.NoError(t, err, "Refresh workspace notification should not error")
	})

	// Create testdata directory if it doesn't exist
	err = os.MkdirAll(filepath.Join("testdata", "snapshots"), 0755)
	require.NoError(t, err, "Failed to create testdata directory")

	// Stop the client
	client.Stop(ctx)
}

// Additional test to verify basic client functionality
func TestClientBasicFunctionality(t *testing.T) {
	// Test client creation
	t.Run("ClientCreation", func(t *testing.T) {
		workspacePath := "/test/workspace/path"
		client := nxlsclient.NewClient(workspacePath, true)

		assert.NotNil(t, client, "Client should not be nil")
		assert.Equal(t, workspacePath, client.NxWorkspacePath, "Workspace path should match")
		assert.NotNil(t, client.Logger, "Logger should be initialized")
	})

	// Test safe client stop when Commander is nil
	t.Run("SafeClientStop", func(t *testing.T) {
		client := nxlsclient.NewClient("/test/path", false)
		assert.Nil(t, client.Commander, "Commander should be nil")

		ctx := context.Background()

		// This should not panic
		client.Stop(ctx)
	})
}

// Helper function to pretty print JSON for better snapshot readability
func prettyPrintJSON(v interface{}) (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

