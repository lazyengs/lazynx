package nxlsclient

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock for ReadWriteCloser testing
type mockReadWriteCloser struct {
	closeCount int
}

func (m *mockReadWriteCloser) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (m *mockReadWriteCloser) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockReadWriteCloser) Close() error {
	m.closeCount++
	return nil
}

func TestReadWriteCloser(t *testing.T) {
	// Create stdin and stdout pipes for testing
	stdoutR, stdoutW := io.Pipe()
	stdinR, stdinW := io.Pipe()
	
	// Create a ReadWriteCloser with our pipes
	rwc := &ReadWriteCloser{
		stdin:  stdinW,
		stdout: stdoutR,
	}
	
	// Test writing
	go func() {
		// Write message to stdin pipe (will be read by our rwc.Read)
		_, err := stdoutW.Write([]byte("test message"))
		assert.NoError(t, err)
		
		// Close pipes after testing
		stdoutW.Close()
	}()
	
	// Test reading
	buf := make([]byte, 12)
	n, err := rwc.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, 12, n)
	assert.Equal(t, "test message", string(buf))
	
	// Test writing
	done := make(chan struct{})
	go func() {
		// Read from the stdinR pipe what our rwc.Write writes
		readBuf := make([]byte, 13)
		n, err := stdinR.Read(readBuf)
		assert.NoError(t, err)
		assert.Equal(t, 13, n)
		assert.Equal(t, "hello, world!", string(readBuf))
		close(done)
	}()
	
	// Write to the rwc
	n, err = rwc.Write([]byte("hello, world!"))
	assert.NoError(t, err)
	assert.Equal(t, 13, n)
	
	// Close the rwc
	err = rwc.Close()
	assert.NoError(t, err)
	
	// Clean up
	stdinR.Close()
}