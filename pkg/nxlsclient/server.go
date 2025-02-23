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

// runSystemReport runs a system report by executing node commands in the server folder.
func (c *Client) runSystemReport(ctx context.Context) error {
	c.Logger.Debugw("System Report:", "serverDir", c.serverDir)
	err := c.runOSCommandInServerFolder(ctx, "node", "-v")
	if err != nil {
		return err
	}
	return c.runOSCommandInServerFolder(ctx, "node", "-p", "process.arch")
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
	c.Logger.Debugw("Stopping nxls")

	// Send requests to stop the NX daemon, shutdown, and exit
	c.Commander.SendStopNxDaemonRequest(ctx)
	c.Commander.SendShutdownRequest(ctx)
	c.Commander.SendExitNotification(ctx)

	// Clean up the server folder
	err := c.cleanUpServerFolder()
	if err != nil {
		return fmt.Errorf("failed to clean up server folder: %w", err)
	}

	// Close the connection if it exists
	if c.conn != nil {
		c.conn.Close()
	}

	// Log the completion of the cleanup process
	c.Logger.Debugw("Cleanup process completed")

	return nil
}

// cleanUpServerFolder removes the temporary server directory.
func (c *Client) cleanUpServerFolder() error {
	err := os.RemoveAll(c.serverDir)
	if err != nil {
		return fmt.Errorf("failed to remove the server directory: %w", err)
	}

	return nil
}
