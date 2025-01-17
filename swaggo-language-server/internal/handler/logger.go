package handler

import (
	"io"
	"log"
)

type logger struct {
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

func NewLogger(writer io.Writer, level LogLevel) *logger {
	return &logger{
		level:          level,
		internalLogger: log.New(writer, "", log.LstdFlags),
	}
}

func (l *logger) Debugf(format string, v ...any) {
	if l.level >= LogDebug {
		l.internalLogger.Printf(format, v...)
	}
}

func (l *logger) Infof(format string, v ...any) {
	if l.level >= LogInfo {
		l.internalLogger.Printf(format, v...)
	}
}

func (l *logger) Warnf(format string, v ...any) {
	if l.level >= LogWarn {
		l.internalLogger.Printf(format, v...)
	}
}

func (l *logger) Error(err error) {
	if l.level >= LogError {
		l.internalLogger.Println(err)
	}
}

func (l *logger) Errorf(format string, v ...any) {
	if l.level >= LogError {
		l.internalLogger.Printf(format, v...)
	}
}
