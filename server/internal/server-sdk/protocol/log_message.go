package protocol

import "github.com/takaaa220/go-swag-ide/server/internal/server-sdk/transport"

type logMessageParams struct {
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

type logger struct {
	maxLogLevel LogLevel
}

func NewLogger(maxLogLevel LogLevel) *logger {
	return &logger{
		maxLogLevel: maxLogLevel,
	}
}

func (l *logger) Log(ctx transport.Context, message string) error {
	if !l.isLog(MessageTypeLog) {
		return nil
	}

	return ctx.Notify("window/logMessage", logMessageParams{
		Type:    MessageTypeLog,
		Message: message,
	})
}

func (l *logger) Info(ctx transport.Context, message string) error {
	if !l.isLog(MessageTypeInfo) {
		return nil
	}

	return ctx.Notify("window/logMessage", logMessageParams{
		Type:    MessageTypeInfo,
		Message: message,
	})
}

func (l *logger) Warn(ctx transport.Context, message string) error {
	if !l.isLog(MessageTypeWarning) {
		return nil
	}

	return ctx.Notify("window/logMessage", logMessageParams{
		Type:    MessageTypeInfo,
		Message: message,
	})
}

func (l *logger) Error(ctx transport.Context, message string) error {
	if !l.isLog(MessageTypeError) {
		return nil
	}

	return ctx.Notify("window/logMessage", logMessageParams{
		Type:    MessageTypeError,
		Message: message,
	})
}

func (l *logger) isLog(logType LogLevel) bool {
	return logType <= l.maxLogLevel
}
