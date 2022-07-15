package logging


import (
	"log"
)

type Level int64

var defaultLogger Loggable = &DefaultLogger {
	defaultLogger: log.Default(),
	Level: INFO,
}

var Logger Loggable = defaultLogger

var fileLogger Loggable

const (
	OFF = iota
	ERROR 
	WARNING
	INFO
	DEBUG
	TRACE
)

func LevelToString(level Level) string {
	mapping := map[Level] string {
		OFF: "OFF",
		ERROR: "ERROR",
		WARNING: "WARNING",
		INFO: "INFO",
		DEBUG: "DEBUG",
		TRACE: "TRACE",
	}
	return mapping[level]
}

type Loggable interface {
	Warn(string)
	Info(string)
	Error(string)
	Debug(string)
	Trace(string)
	SetLevel(Level)
}

type DefaultLogger struct {
	defaultLogger *log.Logger
	Level Level
}

func (l *DefaultLogger) Warn(message string) {
	if l.Level >= WARNING {
		l.defaultLogger.SetPrefix("WARNING:")
		l.defaultLogger.Println(message)
	}
}

func (l *DefaultLogger) Info(message string) {
	if l.Level >= INFO {
		l.defaultLogger.SetPrefix("INFO:")
		l.defaultLogger.Println(message)
	}
}

func (l *DefaultLogger) Error(message string) {
	if l.Level >= ERROR {
		l.defaultLogger.SetPrefix("ERROR:")
		l.defaultLogger.Println(message)
	}
}
func (l *DefaultLogger) Debug(message string) {
	if l.Level >= DEBUG {
		l.defaultLogger.SetPrefix("DEBUG:")
		l.defaultLogger.Println(message)
	}
}
func (l *DefaultLogger) Trace(message string) {
	if l.Level >= TRACE {
		l.defaultLogger.SetPrefix("TRACE:")
		l.defaultLogger.Println(message)
	}
}
func (l *DefaultLogger) SetLevel(level Level) {
	l.Level = level
}