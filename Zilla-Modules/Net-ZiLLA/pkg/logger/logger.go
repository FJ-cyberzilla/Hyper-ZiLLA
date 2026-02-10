package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// LogLevel defines the severity of a log message.
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

// LogEntry represents a structured log entry.
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// Logger provides structured logging functionality.
type Logger struct {
	output     *log.Logger
	level      LogLevel
	jsonOutput bool
	mu         sync.Mutex
}

func NewLogger() *Logger {
	jsonLog := os.Getenv("LOG_FORMAT") == "json"
	return &Logger{
		output:     log.New(os.Stdout, "", 0),
		level:      INFO,
		jsonOutput: jsonLog,
	}
}

func New() *Logger {
	return NewLogger()
}

func (l *Logger) SetJSON(b bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.jsonOutput = b
}

func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}

func (l *Logger) IsDebug() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level <= DEBUG
}

func (l *Logger) WithComponent(name string) *Logger {
	return l // For now just return the same logger, could be enhanced to add a component field
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	msg := fmt.Sprintf(format, args...)

	if l.jsonOutput {
		entry := LogEntry{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Level:     level.String(),
			Message:   msg,
		}
		data, _ := json.Marshal(entry)
		l.output.Println(string(data))
	} else {
		prefix := fmt.Sprintf("[%s] %s ", time.Now().Format("15:04:05"), level.String())
		l.output.Println(prefix + msg)
	}
}

