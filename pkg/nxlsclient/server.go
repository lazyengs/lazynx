package nxlsclient

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"
)

//go:embed server/nxls
var serverfs embed.FS

// unpackServer unpacks the embedded nxls server to a temporary directory.
func (c *Client) unpackServer() error {
	tempDir, err := os.MkdirTemp("", "nxls-server")
	if err != nil {
		return fmt.Errorf("failed to create the temp directory: %w", err)
	}
	c.Logger.Debugw("Created temporary directory", "tempDir", tempDir)

	err = os.CopyFS(tempDir, serverfs)
	if err != nil {
		return fmt.Errorf("failed to copy the server to the temp directory: %w", err)
	}
	c.serverDir = path.Join(tempDir, "server", "nxls")

	return nil
}

// installDependencies installs npm dependencies in the server folder.
func (c *Client) installDependencies(ctx context.Context) error {
	c.Logger.Debugw("Installing dependencies at ", "serverDir", c.serverDir)
	return c.runOSCommandInServerFolder(ctx, "npm", "install")
}

// runOSCommandInServerFolder runs an OS command in the server folder and logs the output.
func (c *Client) runOSCommandInServerFolder(ctx context.Context, name string, args ...string) error {
	c.Logger.Debugw("Running command", "serverDir", c.serverDir, "command", name, "args", args)
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = c.serverDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	if c.isVerbose {
		go func() {
			scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
			for scanner.Scan() {
				c.Logger.Debugw(scanner.Text())
			}
		}()
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("failed to run the command: %w", err)
	}

	return nil
}

// startNxls starts the nxls server and creates the jsonrpc2 connection.
func (c *Client) startNxls(ctx context.Context) (rwc *ReadWriteCloser, err error) {
	serverPath := filepath.Join(c.serverDir, "main.js")

	c.Logger.Debugw("Starting nxls", "workspace", c.NxWorkspacePath, "serverPath", serverPath)

	cmd := exec.CommandContext(ctx, "node", serverPath, "--stdio")
	cmd.Dir = c.NxWorkspacePath

	// Set up process group isolation (prevents signal propagation)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // Put the child in its own process group
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to run the start command: %w", err)
	}

	rwc = &ReadWriteCloser{
		stdin:  stdin,
		stdout: stdout,
	}

	// Start a goroutine to handle the command's completion
	go func() {
		if err := cmd.Wait(); err != nil {
			c.Logger.Errorw("Command exited with error", "error", err)
		}
	}()

	return rwc, nil
}

func (c *Client) stopNxls(ctx context.Context) error {
	// Log the start of the stopping process
	c.Logger.Infow("Stopping nxls server and NX daemon")

	var daemonStoppedWithLSP bool

	// Try LSP commands to stop everything gracefully if Commander is available
	if c.Commander != nil && c.conn != nil {
		c.Logger.Debugw("Attempting to stop NX daemon via LSP protocol")

		// Try to stop the NX daemon via LSP protocol
		err := c.Commander.SendStopNxDaemonRequest(ctx)
		if err != nil {
			c.Logger.Warnw("Failed to stop NX daemon via LSP", "error", err.Error())
		} else {
			c.Logger.Infow("Successfully stopped NX daemon via LSP")
			daemonStoppedWithLSP = true
		}

		c.Logger.Debugw("Sending LSP shutdown request")
		err = c.Commander.SendShutdownRequest(ctx)
		if err != nil {
			c.Logger.Warnw("Failed to send LSP shutdown request", "error", err.Error())
		}

		c.Logger.Debugw("Sending LSP exit notification")
		err = c.Commander.SendExitNotification(ctx)
		if err != nil {
			c.Logger.Warnw("Failed to send LSP exit notification", "error", err.Error())
		}
	} else {
		// Log why we're skipping LSP commands
		if c.Commander == nil {
			c.Logger.Warnw("Commander is nil, skipping LSP requests")
		}
		if c.conn == nil {
			c.Logger.Warnw("Connection is nil, skipping LSP requests")
		}
	}

	// If we couldn't stop the daemon with LSP, try with npx
	if !daemonStoppedWithLSP {
		c.Logger.Infow("Falling back to npx to stop NX daemon")
		err := c.killDaemonWithNpx(ctx)
		if err != nil {
			c.Logger.Errorw("Failed to stop NX daemon with npx", "error", err.Error())
		} else {
			c.Logger.Infow("Successfully stopped NX daemon with npx")
		}
	}

	// Cleanup actions regardless of daemon stop success

	// Clean up the server folder
	c.Logger.Debugw("Cleaning up server folder")
	err := c.cleanUpServerFolder()
	if err != nil {
		c.Logger.Errorw("Failed to clean up server folder", "error", err.Error())
		return fmt.Errorf("failed to clean up server folder: %w", err)
	}

	// Close the connection if it exists, always as the last step
	if c.conn != nil {
		c.Logger.Debugw("Closing LSP connection")
		c.conn.Close()
		c.conn = nil // Set to nil to prevent double-closing
	}

	// Final cleanup status
	if daemonStoppedWithLSP {
		c.Logger.Infow("Cleanup completed successfully (daemon stopped via LSP)")
	} else {
		c.Logger.Infow("Cleanup completed successfully (daemon stopped via npx)")
	}
	return nil
}

func (c *Client) killDaemonWithNpx(ctx context.Context) error {
	c.Logger.Debugw("Attempting to stop NX daemon using npx")

	cmd := exec.CommandContext(ctx, "npx", "nx", "daemon", "--stop")
	cmd.Dir = c.NxWorkspacePath

	// Get stdout and stderr to log the output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start npx nx daemon --stop: %w", err)
	}

	// Read output for logging
	go func() {
		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		for scanner.Scan() {
			c.Logger.Debugw("NX daemon stop output", "message", scanner.Text())
		}
	}()

	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("failed to stop nx daemon with npx: %w", err)
	}

	c.Logger.Debugw("Successfully stopped NX daemon using npx")
	return nil
}

// cleanUpServerFolder removes the temporary server directory.
func (c *Client) cleanUpServerFolder() error {
	// Skip if serverDir is empty
	if c.serverDir == "" {
		c.Logger.Debugw("Server directory not set, skipping cleanup")
		return nil
	}

	// Check if directory exists before attempting to remove
	_, err := os.Stat(c.serverDir)
	if os.IsNotExist(err) {
		c.Logger.Debugw("Server directory doesn't exist, skipping cleanup", "serverDir", c.serverDir)
		return nil
	}

	// Remove the directory
	err = os.RemoveAll(c.serverDir)
	if err != nil {
		return fmt.Errorf("failed to remove the server directory: %w", err)
	}

	c.Logger.Debugw("Server directory removed", "serverDir", c.serverDir)
	return nil
}
