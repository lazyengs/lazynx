package nxlsclient_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/lazyengs/pkg/nxlsclient"
	"github.com/lazyengs/pkg/nxlsclient/commands"
	nxtypes "github.com/lazyengs/pkg/nxlsclient/nx-types"
	"github.com/stretchr/testify/assert"
	"go.lsp.dev/protocol"
)

// Tests that we can create a client correctly
func TestClientCreation(t *testing.T) {
	workspacePath := "/test/workspace/path"
	client := nxlsclient.NewClient(workspacePath, true)

	assert.NotNil(t, client, "Client should not be nil")
	assert.Equal(t, workspacePath, client.NxWorkspacePath, "Workspace path should match")
	assert.NotNil(t, client.Logger, "Logger should be initialized")
}

// Test the integration with Commander
func TestCreateCommander(t *testing.T) {
	// This is a mock test of just the Commander creation
	client := nxlsclient.NewClient("/test/path", false)

	// Since we can't directly create a jsonrpc2.Conn mock that's compatible with the interface,
	// we'll just test that the NewCommander function works with a valid logger

	// Use the unexported method to test commander creation with nil conn
	// This is a partial test - we can't fully test this without exposing internals
	commander := commands.NewCommander(nil, client.Logger)
	assert.NotNil(t, commander, "Commander should be created")
}

// Tests the cleanup process
func TestCleanUp(t *testing.T) {
	// Create a temp directory to simulate server dir
	tempDir, err := os.MkdirTemp("", "nxlsclient-test")
	if err != nil {
		t.Fatal("Failed to create temp dir:", err)
	}

	// We don't actually need the client for this test, as we're just testing
	// directory removal capabilities in general. The client.cleanUpServerFolder method
	// is unexported, so we can't test it directly without modifying the code
	// This test just verifies that directory removal works as expected

	// Create a test file in the temp dir to verify deletion
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatal("Failed to create test file:", err)
	}

	// Verify file exists
	_, err = os.Stat(testFile)
	assert.NoError(t, err, "Test file should exist")

	// Clean up the temp dir manually (simulating what client.cleanUpServerFolder would do)
	err = os.RemoveAll(tempDir)
	assert.NoError(t, err, "Temp dir should be removed without error")

	// Verify directory is gone
	_, err = os.Stat(tempDir)
	assert.True(t, os.IsNotExist(err), "Directory should be removed")
}

// TestClientLifecycleWithNilCommander tests that the client can be stopped
// even when Commander is nil
func TestClientLifecycleWithNilCommander(t *testing.T) {
	// Create a new client
	client := nxlsclient.NewClient("/test/path", false)

	// The Commander will be nil here
	assert.Nil(t, client.Commander, "Commander should be nil")

	// Create context
	ctx := context.Background()

	// This should not panic with the improved code
	client.Stop(ctx)
}

// Integration test - will execute real commands if not skipped
func TestClientLifecycle_Integration(t *testing.T) {
	// Find the nx workspace path - use the project root directory
	workspacePath, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create a new client
	client := nxlsclient.NewClient(workspacePath, false)
	assert.NotNil(t, client, "Client should not be nil")

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
		assert.NotNil(t, initResult, "Initialization result should not be nil")
	case <-time.After(20 * time.Second):
		t.Fatal("Timeout waiting for initialization")
	}

	// Execute basic commands
	if client.Commander != nil {
		_, err := client.Commander.SendWorkspaceRequest(ctx, &commands.WorkspaceRequestParams{Reset: false})
		assert.NoError(t, err, "Workspace request should not error")
	}

	// Stop the client
	client.Stop(ctx)
}

// We would need to create a proper mock for jsonrpc2.Conn if we wanted to fully test
// the client/commander interaction, but for now we're just testing basic functionality

// Mock for client testing
type MockClient struct {
	WorkspacePath string
	ServerDir     string
	Commander     *MockCommander
	IsClosed      bool
}

func NewMockClient(workspacePath string) *MockClient {
	return &MockClient{
		WorkspacePath: workspacePath,
		ServerDir:     "/tmp/mock-server",
		Commander:     NewMockCommander(),
		IsClosed:      false,
	}
}

func (c *MockClient) Start() error {
	return nil
}

func (c *MockClient) Stop() error {
	c.IsClosed = true
	return nil
}

// Mock Commander for testing
type MockCommander struct {
	LastMethod string
	LastParams interface{}
	Results    map[string]interface{}
	Errors     map[string]error
}

func NewMockCommander() *MockCommander {
	return &MockCommander{
		Results: make(map[string]interface{}),
		Errors:  make(map[string]error),
	}
}

func (c *MockCommander) SetResult(method string, result interface{}) {
	c.Results[method] = result
}

func (c *MockCommander) SetError(method string, err error) {
	c.Errors[method] = err
}

func (c *MockCommander) SendWorkspaceRequest(ctx context.Context, params interface{}) (interface{}, error) {
	method := "workspace"
	c.LastMethod = method
	c.LastParams = params

	if err, ok := c.Errors[method]; ok && err != nil {
		return nil, err
	}

	if result, ok := c.Results[method]; ok {
		return result, nil
	}

	// Default mock result
	return map[string]interface{}{
		"workspace": map[string]interface{}{
			"projects": map[string]interface{}{},
		},
	}, nil
}

func (c *MockCommander) SendVersionRequest(ctx context.Context) (*nxtypes.NxVersion, error) {
	method := "nx/version"
	c.LastMethod = method

	if err, ok := c.Errors[method]; ok && err != nil {
		return nil, err
	}

	if result, ok := c.Results[method]; ok {
		if nxVersion, ok := result.(*nxtypes.NxVersion); ok {
			return nxVersion, nil
		}
	}

	// Default mock result
	return &nxtypes.NxVersion{
		Full:  "20.0.0",
		Major: 20,
		Minor: 0,
	}, nil
}

// Test using the mock objects
func TestClientWithMocks(t *testing.T) {
	// Create a mock client
	mockClient := NewMockClient("/test/workspace")

	// Set up mock responses
	version := &nxtypes.NxVersion{
		Full:  "20.0.0",
		Major: 20,
		Minor: 0,
	}
	mockClient.Commander.SetResult("nx/version", version)

	// Set up an error for a different method
	mockClient.Commander.SetError("nonexistent", fmt.Errorf("test error"))

	// Start the client
	err := mockClient.Start()
	assert.NoError(t, err, "Start should not error")

	// Test commander methods
	ctx := context.Background()

	// Test SendVersionRequest with mocked response
	versionResult, err := mockClient.Commander.SendVersionRequest(ctx)
	assert.NoError(t, err, "SendVersionRequest should not error")
	assert.Equal(t, "20.0.0", versionResult.Full, "Version should match mock")

	// Test SendWorkspaceRequest with default mock response
	workspace, err := mockClient.Commander.SendWorkspaceRequest(ctx, nil)
	assert.NoError(t, err, "SendWorkspaceRequest should not error")
	assert.NotNil(t, workspace, "Workspace should not be nil")

	// Stop the client
	err = mockClient.Stop()
	assert.NoError(t, err, "Stop should not error")
	assert.True(t, mockClient.IsClosed, "Client should be marked as closed")
}
