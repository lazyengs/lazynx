package nxlsclient

import (
	"encoding/json"
	"sync"
	"sync/atomic"
)

// NotificationHandler is a function that handles a notification.
// It takes a notification method and parameters as input.
type NotificationHandler func(method string, params json.RawMessage) error

// handlerEntry represents a single registered handler with a unique ID.
type handlerEntry struct {
	id      uint64
	handler NotificationHandler
}

// Disposable provides a way to dispose/unregister a notification handler.
type Disposable struct {
	id       uint64
	method   string
	listener *notificationListener
}

// Dispose unregisters the handler associated with this disposable.
func (d *Disposable) Dispose() {
	if d.listener != nil {
		d.listener.unregisterHandlerByID(d.method, d.id)
	}
}

// NotificationListener manages notification handlers for different notification methods.
type notificationListener struct {
	mu        sync.RWMutex
	handlers  map[string][]handlerEntry
	idCounter atomic.Uint64
}

// NewNotificationListener creates a new NotificationListener instance.
func newNotificationListener() *notificationListener {
	return &notificationListener{
		handlers: make(map[string][]handlerEntry),
	}
}

// RegisterHandler registers a handler for a specific notification method.
// Returns a Disposable that can be used to unregister the handler.
func (l *notificationListener) registerHandler(method string, handler NotificationHandler) *Disposable {
	if handler == nil {
		return &Disposable{listener: nil}
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Generate a unique ID for this handler
	handlerID := l.idCounter.Add(1)

	// Create the entry
	entry := handlerEntry{
		id:      handlerID,
		handler: handler,
	}

	// Add the handler to the map
	l.handlers[method] = append(l.handlers[method], entry)

	// Return a disposable to unregister the handler
	return &Disposable{
		id:       handlerID,
		method:   method,
		listener: l,
	}
}

// unregisterHandlerByID removes a handler with the specified ID for a method.
func (l *notificationListener) unregisterHandlerByID(method string, id uint64) {
	l.mu.Lock()
	defer l.mu.Unlock()

	handlers, ok := l.handlers[method]
	if !ok {
		return
	}

	// Find and remove the handler with matching ID
	for i, entry := range handlers {
		if entry.id == id {
			// Remove this handler by appending slices before and after it
			l.handlers[method] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	// If no handlers left for this method, remove the method entry
	if len(l.handlers[method]) == 0 {
		delete(l.handlers, method)
	}
}

// NotifyAll calls all registered handlers for a specific notification method.
func (l *notificationListener) notifyAll(method string, params json.RawMessage) {
	l.mu.RLock()
	entries, ok := l.handlers[method]
	l.mu.RUnlock()

	if !ok {
		return
	}

	// Make a copy of handlers to avoid holding the lock during execution
	entriesCopy := make([]handlerEntry, len(entries))
	copy(entriesCopy, entries)

	// Call each handler
	for _, entry := range entriesCopy {
		// Errors from handlers are intentionally ignored to ensure all handlers are called
		_ = entry.handler(method, params)
	}
}

// ClearHandlers removes all handlers for all methods.
func (l *notificationListener) clearHandlers() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.handlers = make(map[string][]handlerEntry)
}

// HasHandlers checks if there are any handlers registered for a specific method.
func (l *notificationListener) hasHandlers(method string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	handlers, ok := l.handlers[method]
	return ok && len(handlers) > 0
}
