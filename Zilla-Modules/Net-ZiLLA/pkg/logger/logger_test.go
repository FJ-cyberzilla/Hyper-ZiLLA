package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestLogLevel_String(t *testing.T) {
	if DEBUG.String() != "DEBUG" {
		t.Errorf("expected DEBUG, got %s", DEBUG.String())
	}
	if FATAL.String() != "FATAL" {
		t.Errorf("expected FATAL, got %s", FATAL.String())
	}
}

func TestLogger_Levels(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{
		output:     log.New(&buf, "", 0),
		level:      INFO,
		jsonOutput: false,
	}

	l.Debug("debug message")
	if buf.Len() > 0 {
		t.Errorf("expected no debug message, got %s", buf.String())
	}

	l.Info("info message")
	if !strings.Contains(buf.String(), "INFO") || !strings.Contains(buf.String(), "info message") {
		t.Errorf("expected info message, got %s", buf.String())
	}
	buf.Reset()

	l.SetLevel(DEBUG)
	l.Debug("debug message")
	if !strings.Contains(buf.String(), "DEBUG") || !strings.Contains(buf.String(), "debug message") {
		t.Errorf("expected debug message, got %s", buf.String())
	}
}

func TestLogger_JSON(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{
		output:     log.New(&buf, "", 0),
		level:      INFO,
		jsonOutput: true,
	}

	l.Info("json message")
	if !strings.Contains(buf.String(), "\"level\":\"INFO\"") || !strings.Contains(buf.String(), "\"message\":\"json message\"") {
		t.Errorf("expected json log, got %s", buf.String())
	}
}

func TestLogger_SetJSON(t *testing.T) {
	l := NewLogger()
	l.SetJSON(true)
	if !l.jsonOutput {
		t.Errorf("expected jsonOutput to be true")
	}
	l.SetJSON(false)
	if l.jsonOutput {
		t.Errorf("expected jsonOutput to be false")
	}
}
