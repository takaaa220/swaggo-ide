package protocol

// LogMessageParams defines the parameters for the `window/logMessage` notification.
type LogMessageParams struct {
	Type    LogMessageType `json:"type"`
	Message string         `json:"message"`
}

// MessageType represents the type of message to log.
type LogMessageType int

const (
	MessageTypeError   LogMessageType = 1
	MessageTypeWarning LogMessageType = 2
	MessageTypeInfo    LogMessageType = 3
	MessageTypeLog     LogMessageType = 4
)
