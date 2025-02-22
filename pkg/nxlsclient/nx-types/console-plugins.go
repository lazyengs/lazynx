package nxtypes

type StartupMessageType string

const (
	StartupMessageTypeWarning StartupMessageType = "warning"
	StartupMessageTypeError   StartupMessageType = "error"
)

type StartupMessageDefinition struct {
	Message string             `json:"message"`
	Type    StartupMessageType `json:"type"`
}
