package nxlsclient

import (
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNotificationListener(t *testing.T) {
	t.Run("registerHandler", func(t *testing.T) {
		listener := newNotificationListener()

		// Register a handler
		called := false
		disposable := listener.registerHandler("test/method", func(method string, params json.RawMessage) error {
			called = true
			assert.Equal(t, "test/method", method)
			return nil
		})

		// Check if handler is registered
		assert.True(t, listener.hasHandlers("test/method"), "Handler should be registered")

		// Notify and check if handler was called
		rawJSON := json.RawMessage(`{"key": "value"}`)
		listener.notifyAll("test/method", rawJSON)
		assert.True(t, called, "Handler should have been called")

		// Test unregistration
		disposable.Dispose()
		assert.False(t, listener.hasHandlers("test/method"))

		// Reset the called flag
		called = false

		// Notify again, handler should not be called
		listener.notifyAll("test/method", rawJSON)
		assert.False(t, called)
	})

	t.Run("MultipleHandlers", func(t *testing.T) {
		listener := newNotificationListener()

		// Track call counts
		var mu sync.Mutex
		callCount1 := 0
		callCount2 := 0

		// Register multiple handlers for the same method
		disposable1 := listener.registerHandler("test/method", func(method string, params json.RawMessage) error {
			mu.Lock()
			defer mu.Unlock()
			callCount1++
			return nil
		})

		disposable2 := listener.registerHandler("test/method", func(method string, params json.RawMessage) error {
			mu.Lock()
			defer mu.Unlock()
			callCount2++
			return nil
		})

		// Notify and check if both handlers were called
		rawJSON := json.RawMessage(`{"key": "value"}`)
		listener.notifyAll("test/method", rawJSON)

		// Short sleep to allow for any async processing
		time.Sleep(10 * time.Millisecond)

		mu.Lock()
		assert.Equal(t, 1, callCount1, "First handler should have been called once")
		assert.Equal(t, 1, callCount2, "Second handler should have been called once")
		mu.Unlock()

		// Unregister one handler and try again
		disposable1.Dispose()
		listener.notifyAll("test/method", rawJSON)

		// Short sleep to allow for any async processing
		time.Sleep(10 * time.Millisecond)

		mu.Lock()
		assert.Equal(t, 1, callCount1, "First handler should not have been called again")
		assert.Equal(t, 2, callCount2, "Second handler should have been called again")
		mu.Unlock()

		// Clean up
		disposable2.Dispose()
	})

	t.Run("clearHandlers", func(t *testing.T) {
		listener := newNotificationListener()

		// Register handlers for different methods
		listener.registerHandler("method1", func(method string, params json.RawMessage) error {
			return nil
		})

		listener.registerHandler("method2", func(method string, params json.RawMessage) error {
			return nil
		})

		assert.True(t, listener.hasHandlers("method1"))
		assert.True(t, listener.hasHandlers("method2"))

		// Clear all handlers
		listener.clearHandlers()

		assert.False(t, listener.hasHandlers("method1"))
		assert.False(t, listener.hasHandlers("method2"))
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		listener := newNotificationListener()

		// Create a WaitGroup to wait for all goroutines to finish
		var wg sync.WaitGroup

		// Number of concurrent registrations
		numConcurrent := 100
		wg.Add(numConcurrent * 2) // For both registration and notification

		// Track all disposables to prevent them from being garbage collected
		disposables := make([]*Disposable, 0, numConcurrent)
		var disposablesMu sync.Mutex

		// Concurrently register handlers
		for i := 0; i < numConcurrent; i++ {
			go func(index int) {
				defer wg.Done()
				d := listener.registerHandler("concurrent/test", func(method string, params json.RawMessage) error {
					return nil
				})

				disposablesMu.Lock()
				disposables = append(disposables, d)
				disposablesMu.Unlock()
			}(i)
		}

		// Concurrently notify
		for i := 0; i < numConcurrent; i++ {
			go func(index int) {
				defer wg.Done()
				rawJSON := json.RawMessage(`{"index": ` + string([]byte{byte('0' + index%10)}) + `}`)
				listener.notifyAll("concurrent/test", rawJSON)
			}(i)
		}

		// Wait for all goroutines to complete
		wg.Wait()

		// Verify handlers were registered
		assert.True(t, listener.hasHandlers("concurrent/test"))
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		listener := newNotificationListener()

		// Register handlers that return errors, but expect the notification to continue
		var callCount int

		listener.registerHandler("error/test", func(method string, params json.RawMessage) error {
			callCount++
			return errors.New("handler error")
		})

		listener.registerHandler("error/test", func(method string, params json.RawMessage) error {
			callCount++
			return nil
		})

		// Notify should call both handlers regardless of errors
		rawJSON := json.RawMessage(`{"key": "value"}`)
		listener.notifyAll("error/test", rawJSON)

		assert.Equal(t, 2, callCount)
	})

	t.Run("NilHandler", func(t *testing.T) {
		listener := newNotificationListener()

		// Register a nil handler (should be a no-op)
		disposable := listener.registerHandler("test/method", nil)

		// Should not register the nil handler
		assert.False(t, listener.hasHandlers("test/method"))

		// Dispose should be safe to call
		disposable.Dispose()
	})

	t.Run("DisposeTwice", func(t *testing.T) {
		// Test that disposing twice doesn't panic
		listener := newNotificationListener()

		disposable := listener.registerHandler("test/method", func(method string, params json.RawMessage) error {
			return nil
		})

		// First dispose should work
		disposable.Dispose()
		assert.False(t, listener.hasHandlers("test/method"))

		// Second dispose should be a no-op and not panic
		disposable.Dispose()
	})
}
