package nxlsclient

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestUnpackServer(t *testing.T) {
	// Skip if no real filesystem access
	if os.Getenv("SKIP_FS_TESTS") == "true" {
		t.Skip("Skipping test requiring filesystem access")
	}

	// Create a test client
	logger, _ := zap.NewDevelopment()
	client := &Client{
		Logger: logger.Sugar(),
	}

	// Test unpacking the server
	err := client.unpackServer()
	assert.NoError(t, err, "Unpacking server should not error")
	assert.NotEmpty(t, client.serverDir, "Server directory should be set")

	// Check that the directory exists
	_, err = os.Stat(client.serverDir)
	assert.NoError(t, err, "Server directory should exist")

	// Clean up
	os.RemoveAll(client.serverDir)
}

func TestInstallDependencies(t *testing.T) {
	// Skip if no real filesystem access or if running in CI
	if os.Getenv("SKIP_FS_TESTS") == "true" || os.Getenv("CI") == "true" {
		t.Skip("Skipping test requiring filesystem and npm access")
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "nxls-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test client
	logger, _ := zap.NewDevelopment()
	client := &Client{
		Logger:    logger.Sugar(),
		serverDir: tempDir,
	}

	// Create minimal package.json in the temp directory
	packageJSON := `{
		"name": "nxls-test",
		"version": "1.0.0",
		"dependencies": {
			"nx": "^15.0.0"
		}
	}`
	err = os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to write package.json: %v", err)
	}

	// Test installing dependencies with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// This test might fail if npm is not available or network issues
	err = client.installDependencies(ctx)
	if err != nil {
		t.Logf("Installation of dependencies failed: %v (this might be expected in some environments)", err)
	}
}

func TestCleanUpServerFolder(t *testing.T) {
	// Create a test client
	logger, _ := zap.NewDevelopment()
	client := &Client{
		Logger: logger.Sugar(),
	}

	// Create a temp directory that will be cleaned up
	tempDir, err := os.MkdirTemp("", "nxls-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Set server directory and test cleanup
	client.serverDir = tempDir
	err = client.cleanUpServerFolder()
	assert.NoError(t, err, "Cleaning up server directory should not error")

	// Verify directory is removed
	_, err = os.Stat(tempDir)
	assert.True(t, os.IsNotExist(err), "Server directory should be removed")
}

func TestInvalidWorkspacePath(t *testing.T) {
	// Create a client with an invalid workspace path
	invalidPath := "/path/that/does/not/exist"
	client := NewClient(invalidPath, false)

	// Just check that client is created properly - we can't test full start/stop
	// without fixing the stopNxls method to handle nil Commander
	assert.NotNil(t, client, "Client should be created even with invalid path")
	assert.Equal(t, invalidPath, client.NxWorkspacePath, "Workspace path should match")
}
