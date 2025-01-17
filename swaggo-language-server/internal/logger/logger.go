package logger

import (
	"io"
	"log"
)

var l *lspLogger

type lspLogger struct {
	level          LogLevel
	internalLogger *log.Logger
}

type LogLevel int

const (
	LogDebug LogLevel = iota
	LogInfo
	LogWarn
	LogError
)

func Setup(writer io.Writer, level LogLevel) {
	l = &lspLogger{
		level:          level,
		internalLogger: log.New(writer, "", log.LstdFlags),
	}
}

func Debugf(format string, v ...any) {
	if l.level <= LogDebug {
		l.internalLogger.Printf(format, v...)
	}
}

func Infof(format string, v ...any) {
	if l.level <= LogInfo {
		l.internalLogger.Printf(format, v...)
	}
}

func Warnf(format string, v ...any) {
	if l.level <= LogWarn {
		l.internalLogger.Printf(format, v...)
	}
}

func Error(err error) {
	if l.level <= LogError {
		l.internalLogger.Println(err)
	}
}

func Errorf(format string, v ...any) {
	if l.level <= LogError {
		l.internalLogger.Printf(format, v...)
	}
}
