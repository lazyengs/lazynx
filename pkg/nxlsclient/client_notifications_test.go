package nxlsclient

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/sourcegraph/jsonrpc2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNotificationHandling tests the notification handling capabilities of the client
func TestNotificationHandling(t *testing.T) {
	// Skip in CI
	if testing.Short() {
		t.Skip("Skipping notification test in short mode")
	}

	// Create a new client
	client := NewClient("/test/path", true)
	require.NotNil(t, client, "Client should not be nil")

	// Test registering notification handlers
	t.Run("RegisterNotificationHandler", func(t *testing.T) {
		// Register a handler for window/logMessage
		var receivedMessage string
		var wg sync.WaitGroup
		wg.Add(1)

		// Use the typed notification handler helper
		disposable := client.OnNotification(
			WindowLogMessageMethod,
			TypedNotificationHandler(
				func(method string, params *WindowLogMessage) error {
					receivedMessage = params.Message
					wg.Done()
					return nil
				},
			),
		)

		// Manually trigger the notification
		params := json.RawMessage(`{"message": "Test message", "type": 1}`)
		// Use handleServerRequest directly to simulate receiving a notification
		req := &jsonrpc2.Request{
			Method: WindowLogMessageMethod,
			Params: &params,
			Notif:  true,
		}
		client.handleServerRequest(context.Background(), nil, req)

		// Wait for the notification to be processed
		waitDone := make(chan struct{})
		go func() {
			wg.Wait()
			close(waitDone)
		}()

		select {
		case <-waitDone:
			// Success
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for notification handler to be called")
		}

		// Check that the handler was called with the correct message
		assert.Equal(t, "Test message", receivedMessage)

		// Test unregistering the handler
		disposable.Dispose()

		// Reset and try again, should not be called
		receivedMessage = ""
		wg.Add(1)
		// Use handleServerRequest directly to simulate receiving a notification
		req = &jsonrpc2.Request{
			Method: WindowLogMessageMethod,
			Params: &params,
			Notif:  true,
		}
		client.handleServerRequest(context.Background(), nil, req)

		// Since we unregistered, the handler should not be called
		timeoutCh := time.After(500 * time.Millisecond)
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			t.Fatal("Handler should not have been called after unregistering")
		case <-timeoutCh:
			// This is expected - the handler wasn't called
			wg.Done() // Prevent WaitGroup deadlock
		}

		assert.Empty(t, receivedMessage, "Message should not have been updated after unregistering handler")
	})

	// Test multiple notification handlers for the same method
	t.Run("MultipleHandlers", func(t *testing.T) {
		count1, count2 := 0, 0
		var mu sync.Mutex

		// Register two handlers for the same method
		disposable1 := client.OnNotification(NxRefreshWorkspaceMethod, func(method string, params json.RawMessage) error {
			mu.Lock()
			defer mu.Unlock()
			count1++
			return nil
		})

		disposable2 := client.OnNotification(NxRefreshWorkspaceMethod, func(method string, params json.RawMessage) error {
			mu.Lock()
			defer mu.Unlock()
			count2++
			return nil
		})

		// Manually trigger the notification
		params := json.RawMessage(`{}`) // empty params
		// Use handleServerRequest directly to simulate receiving a notification
		req := &jsonrpc2.Request{
			Method: NxRefreshWorkspaceMethod,
			Params: &params,
			Notif:  true,
		}
		client.handleServerRequest(context.Background(), nil, req)

		// Short sleep to allow the async handlers to execute
		time.Sleep(100 * time.Millisecond)

		mu.Lock()
		assert.Equal(t, 1, count1, "First handler should have been called once")
		assert.Equal(t, 1, count2, "Second handler should have been called once")
		mu.Unlock()

		// Unregister one handler and try again
		disposable1.Dispose()
		// Use handleServerRequest directly to simulate receiving a notification
		req = &jsonrpc2.Request{
			Method: NxRefreshWorkspaceMethod,
			Params: &params,
			Notif:  true,
		}
		client.handleServerRequest(context.Background(), nil, req)

		// Short sleep to allow the async handlers to execute
		time.Sleep(100 * time.Millisecond)

		mu.Lock()
		assert.Equal(t, 1, count1, "First handler should not have been called again")
		assert.Equal(t, 2, count2, "Second handler should have been called again")
		mu.Unlock()

		// Clean up
		disposable2.Dispose()
	})

	// Test client cleanup
	t.Run("ClientCleanup", func(t *testing.T) {
		// Register a handler
		handlerCalled := false
		disposable := client.OnNotification("test/method", func(method string, params json.RawMessage) error {
			handlerCalled = true
			return nil
		})

		// Verify the handler is registered
		assert.True(t, client.notificationListener.hasHandlers("test/method"))
		assert.NotNil(t, disposable)

		// Stop the client
		ctx := context.Background()
		client.Stop(ctx)

		// Verify all handlers were cleared
		assert.False(t, client.notificationListener.hasHandlers("test/method"))

		// Manually trigger the notification (should not call the handler)
		handlerCalled = false
		params := json.RawMessage(`{}`)
		// Use handleServerRequest directly to simulate receiving a notification
		req := &jsonrpc2.Request{
			Method: "test/method",
			Params: &params,
			Notif:  true,
		}
		client.handleServerRequest(context.Background(), nil, req)

		// Verify the handler was not called
		assert.False(t, handlerCalled)

		// Disposing after client stop should be safe
		disposable.Dispose()
	})
}

// TestE2ENotifications is a test that attempts to connect to a real server
// and verify notifications work in practice. This is skipped by default.
func TestE2ENotifications(t *testing.T) {
	// Skip in CI or short test mode
	if testing.Short() {
		t.Skip("Skipping E2E notification test in short mode")
	}

	// Always skip this test as it's just an example
	t.Skip("This is an example test that requires a real NX workspace")
}
