package nxlsclient

import (
	"bufio"
	"context"
	"embed"
	"errors"
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
		c.Logger.Errorf("Failed to create temp directory: %s", err.Error())
		return errors.New("Failed to create the temp directory")
	}
	c.Logger.Debugw("Created temporary directory", "tempDir", tempDir)

	err = os.CopyFS(tempDir, serverfs)
	if err != nil {
		c.Logger.Errorf("Failed to copy the server to the temp directory: %s", err.Error())
		return errors.New("Failed to copy the server to the temp directory")

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
		c.Logger.Fatalf("Failed to get stdout pipe: %s", err.Error())
		return errors.New("Failed to get stdout pipe")
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		c.Logger.Errorf("Failed to get stderr pipe: %s", err.Error())
		return errors.New("Failed to get stderr pipe")
	}

	if err := cmd.Start(); err != nil {
		c.Logger.Errorf("Failed to start command: %s", err.Error())
		return errors.New("Failed to start command")
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
		c.Logger.Errorf("Failed to run the command %s: %s", name, err.Error())
		return errors.New("Failed to run the command")
	}

	return nil
}

// startNxls starts the nxls server and creates the jsonrpc2 connection.
func (c *Client) startNxls(ctx context.Context) (rwc *ReadWriteCloser, err error) {
	serverPath := filepath.Join(c.serverDir, "main.js")

	c.Logger.Debugw("Starting nxls", "workspace", c.nxWorkspacePath, "serverPath", serverPath)

	cmd := exec.CommandContext(ctx, "node", serverPath, "--stdio")
	cmd.Dir = c.nxWorkspacePath

	stdin, err := cmd.StdinPipe()
	if err != nil {
		c.Logger.Fatalf("failed to create stdin pipe: %s", err.Error())
		return nil, errors.New("Failed to get stdin pipe")
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c.Logger.Fatalf("Failed to get stdout pipe: %s", err.Error())
		return nil, errors.New("Failed to get stdout pipe")
	}

	if err := cmd.Start(); err != nil {
		c.Logger.Errorf("Failed to start command: %s", err.Error())
		return nil, errors.New("Failed to start command")
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
