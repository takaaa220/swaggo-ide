package protocol

type LogMessageParams struct {
	Type    LogLevel `json:"type"`
	Message string   `json:"message"`
}

type LogLevel int

const (
	MessageTypeError   LogLevel = 1
	MessageTypeWarning LogLevel = 2
	MessageTypeInfo    LogLevel = 3
	MessageTypeLog     LogLevel = 4
)

func (t LogLevel) String() string {
	switch t {
	case MessageTypeError:
		return "error"
	case MessageTypeWarning:
		return "warning"
	case MessageTypeInfo:
		return "info"
	case MessageTypeLog:
		return "log"
	default:
		return "unknown"
	}
}
