package nxlsclient

import (
	"io"
)

// ReadWriteCloser combines stdin and stdout of a command into a single io.ReadWriteCloser
type ReadWriteCloser struct {
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

func (c *ReadWriteCloser) Read(p []byte) (n int, err error) {
	return c.stdout.Read(p)
}

func (c *ReadWriteCloser) Write(p []byte) (n int, err error) {
	return c.stdin.Write(p)
}

func (c *ReadWriteCloser) Close() error {
	err1 := c.stdin.Close()
	err2 := c.stdout.Close()
	if err1 != nil {
		return err1
	}
	return err2
}
