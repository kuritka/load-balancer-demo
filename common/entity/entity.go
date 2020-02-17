package entity

import "time"

const (
	LogLevelInfo  LogLevel = "INFO"
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelError LogLevel = "ERROR"
)

type LogLevel string

type LogEntry struct {
	Level       LogLevel
	Timestamp   time.Time
	Source      string
	Destination string
	Message     string
}
