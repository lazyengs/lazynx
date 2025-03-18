package nxlsclient

import (
	"encoding/json"
	"fmt"

	"github.com/lazyengs/lazynx/pkg/nxlsclient/commands"
)

const (
	// Notification method constants
	WindowLogMessageMethod          = "window/logMessage"
	NxRefreshWorkspaceMethod        = commands.RefreshWorkspaceNotificationMethod
	NxRefreshWorkspaceStartedMethod = commands.RefreshWorkspaceStartedNotificationMethod
	NxChangeWorkspaceMethod         = commands.ChangeWorkspaceNotificationMethod
)

// WindowLogMessage represents a window/logMessage notification from the server
type WindowLogMessage struct {
	Message string `json:"message"`
	Type    int8   `json:"type"`
}

// ParseNotification is a helper function to parse JSON notification parameters into a typed struct
func ParseNotification[T any](params json.RawMessage) (*T, error) {
	if params == nil {
		return nil, fmt.Errorf("notification params are nil")
	}

	var result T
	if err := json.Unmarshal(params, &result); err != nil {
		return nil, fmt.Errorf("failed to parse notification: %w", err)
	}

	return &result, nil
}

// Helper function to create a strongly typed notification handler
func TypedNotificationHandler[T any](handler func(method string, params *T) error) NotificationHandler {
	return func(method string, params json.RawMessage) error {
		parsed, err := ParseNotification[T](params)
		if err != nil {
			return err
		}

		return handler(method, parsed)
	}
}

// Helper constants for notification message types
const (
	LogError   int8 = 1
	LogWarning int8 = 2
	LogInfo    int8 = 3
	LogDebug   int8 = 4
)
