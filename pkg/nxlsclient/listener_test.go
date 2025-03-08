package nxlsclient_test

import (
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/lazyengs/pkg/nxlsclient"
	"github.com/stretchr/testify/assert"
)

func TestNotificationListener(t *testing.T) {
	t.Run("RegisterHandler", func(t *testing.T) {
		listener := nxlsclient.NewNotificationListener()
		
		// Register a handler
		called := false
		disposable := listener.RegisterHandler("test/method", func(method string, params json.RawMessage) error {
			called = true
			assert.Equal(t, "test/method", method)
			return nil
		})
		
		// Check if handler is registered
		assert.True(t, listener.HasHandlers("test/method"), "Handler should be registered")
		
		// Notify and check if handler was called
		rawJSON := json.RawMessage(`{"key": "value"}`)
		listener.NotifyAll("test/method", rawJSON)
		assert.True(t, called, "Handler should have been called")
		
		// Test unregistration
		disposable.Dispose()
		assert.False(t, listener.HasHandlers("test/method"))
		
		// Reset the called flag
		called = false
		
		// Notify again, handler should not be called
		listener.NotifyAll("test/method", rawJSON)
		assert.False(t, called)
	})
	
	t.Run("MultipleHandlers", func(t *testing.T) {
		listener := nxlsclient.NewNotificationListener()
		
		// Track call counts
		var mu sync.Mutex
		callCount1 := 0
		callCount2 := 0
		
		// Register multiple handlers for the same method
		disposable1 := listener.RegisterHandler("test/method", func(method string, params json.RawMessage) error {
			mu.Lock()
			defer mu.Unlock()
			callCount1++
			return nil
		})
		
		disposable2 := listener.RegisterHandler("test/method", func(method string, params json.RawMessage) error {
			mu.Lock()
			defer mu.Unlock()
			callCount2++
			return nil
		})
		
		// Notify and check if both handlers were called
		rawJSON := json.RawMessage(`{"key": "value"}`)
		listener.NotifyAll("test/method", rawJSON)
		
		// Short sleep to allow for any async processing
		time.Sleep(10 * time.Millisecond)
		
		mu.Lock()
		assert.Equal(t, 1, callCount1, "First handler should have been called once")
		assert.Equal(t, 1, callCount2, "Second handler should have been called once")
		mu.Unlock()
		
		// Unregister one handler and try again
		disposable1.Dispose()
		listener.NotifyAll("test/method", rawJSON)
		
		// Short sleep to allow for any async processing
		time.Sleep(10 * time.Millisecond)
		
		mu.Lock()
		assert.Equal(t, 1, callCount1, "First handler should not have been called again")
		assert.Equal(t, 2, callCount2, "Second handler should have been called again")
		mu.Unlock()
		
		// Clean up
		disposable2.Dispose()
	})
	
	t.Run("ClearHandlers", func(t *testing.T) {
		listener := nxlsclient.NewNotificationListener()
		
		// Register handlers for different methods
		listener.RegisterHandler("method1", func(method string, params json.RawMessage) error {
			return nil
		})
		
		listener.RegisterHandler("method2", func(method string, params json.RawMessage) error {
			return nil
		})
		
		assert.True(t, listener.HasHandlers("method1"))
		assert.True(t, listener.HasHandlers("method2"))
		
		// Clear all handlers
		listener.ClearHandlers()
		
		assert.False(t, listener.HasHandlers("method1"))
		assert.False(t, listener.HasHandlers("method2"))
	})
	
	t.Run("ConcurrentAccess", func(t *testing.T) {
		listener := nxlsclient.NewNotificationListener()
		
		// Create a WaitGroup to wait for all goroutines to finish
		var wg sync.WaitGroup
		
		// Number of concurrent registrations
		numConcurrent := 100
		wg.Add(numConcurrent * 2) // For both registration and notification
		
		// Track all disposables to prevent them from being garbage collected
		disposables := make([]*nxlsclient.Disposable, 0, numConcurrent)
		var disposablesMu sync.Mutex
		
		// Concurrently register handlers
		for i := 0; i < numConcurrent; i++ {
			go func(index int) {
				defer wg.Done()
				d := listener.RegisterHandler("concurrent/test", func(method string, params json.RawMessage) error {
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
				listener.NotifyAll("concurrent/test", rawJSON)
			}(i)
		}
		
		// Wait for all goroutines to complete
		wg.Wait()
		
		// Verify handlers were registered
		assert.True(t, listener.HasHandlers("concurrent/test"))
	})
	
	t.Run("ErrorHandling", func(t *testing.T) {
		listener := nxlsclient.NewNotificationListener()
		
		// Register handlers that return errors, but expect the notification to continue
		var callCount int
		
		listener.RegisterHandler("error/test", func(method string, params json.RawMessage) error {
			callCount++
			return errors.New("handler error")
		})
		
		listener.RegisterHandler("error/test", func(method string, params json.RawMessage) error {
			callCount++
			return nil
		})
		
		// Notify should call both handlers regardless of errors
		rawJSON := json.RawMessage(`{"key": "value"}`)
		listener.NotifyAll("error/test", rawJSON)
		
		assert.Equal(t, 2, callCount)
	})
	
	t.Run("NilHandler", func(t *testing.T) {
		listener := nxlsclient.NewNotificationListener()
		
		// Register a nil handler (should be a no-op)
		disposable := listener.RegisterHandler("test/method", nil)
		
		// Should not register the nil handler
		assert.False(t, listener.HasHandlers("test/method"))
		
		// Dispose should be safe to call
		disposable.Dispose()
	})
	
	t.Run("DisposeTwice", func(t *testing.T) {
		// Test that disposing twice doesn't panic
		listener := nxlsclient.NewNotificationListener()
		
		disposable := listener.RegisterHandler("test/method", func(method string, params json.RawMessage) error {
			return nil
		})
		
		// First dispose should work
		disposable.Dispose()
		assert.False(t, listener.HasHandlers("test/method"))
		
		// Second dispose should be a no-op and not panic
		disposable.Dispose()
	})
}