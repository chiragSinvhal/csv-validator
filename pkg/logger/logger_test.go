package logger

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer

	// Save original loggers
	origDebug := debugLogger
	origInfo := infoLogger
	origWarn := warnLogger
	origError := errorLogger

	// Replace with test loggers
	debugLogger = log.New(&buf, "[DEBUG] ", log.Ldate|log.Ltime)
	infoLogger = log.New(&buf, "[INFO] ", log.Ldate|log.Ltime)
	warnLogger = log.New(&buf, "[WARN] ", log.Ldate|log.Ltime)
	errorLogger = log.New(&buf, "[ERROR] ", log.Ldate|log.Ltime)

	f()

	// Restore
	debugLogger = origDebug
	infoLogger = origInfo
	warnLogger = origWarn
	errorLogger = origError

	return buf.String()
}

func TestInit(t *testing.T) {
	tests := []struct {
		level    string
		expected LogLevel
	}{
		{"debug", DEBUG},
		{"DEBUG", DEBUG},
		{"info", INFO},
		{"INFO", INFO},
		{"warn", WARN},
		{"warning", WARN},
		{"error", ERROR},
		{"ERROR", ERROR},
		{"invalid", INFO}, // defaults to INFO
		{"", INFO},
	}

	for _, tt := range tests {
		Init(tt.level)
		assert.Equal(t, tt.expected, currentLevel, "Failed for level: %s", tt.level)
	}
}

func TestDebug(t *testing.T) {
	Init("debug")
	output := captureOutput(func() {
		Debug("test debug message")
	})
	assert.Contains(t, output, "test debug message")
	assert.Contains(t, output, "[DEBUG]")
}

func TestDebug_NotLogged(t *testing.T) {
	Init("info")
	output := captureOutput(func() {
		Debug("should not appear")
	})
	assert.Empty(t, output)
}

func TestInfo(t *testing.T) {
	Init("info")
	output := captureOutput(func() {
		Info("test info message")
	})
	assert.Contains(t, output, "test info message")
	assert.Contains(t, output, "[INFO]")
}

func TestWarn(t *testing.T) {
	Init("warn")
	output := captureOutput(func() {
		Warn("test warning")
	})
	assert.Contains(t, output, "test warning")
	assert.Contains(t, output, "[WARN]")
}

func TestError(t *testing.T) {
	Init("error")

	// Redirect stderr to capture error logs
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	errorLogger = log.New(w, "[ERROR] ", log.Ldate|log.Ltime)

	Error("test error")

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	assert.Contains(t, output, "test error")
	assert.Contains(t, output, "[ERROR]")
}

func TestLogLevelFiltering(t *testing.T) {
	Init("warn")

	output := captureOutput(func() {
		Debug("debug msg")
		Info("info msg")
		Warn("warn msg")
		Error("error msg")
	})

	// Only warn and error should appear
	assert.NotContains(t, output, "debug msg")
	assert.NotContains(t, output, "info msg")
	assert.Contains(t, output, "warn msg")
}

func TestDebugf(t *testing.T) {
	Init("debug")
	output := captureOutput(func() {
		Debugf("test %s with %d", "format", 123)
	})
	assert.Contains(t, output, "test format with 123")
}

func TestInfof(t *testing.T) {
	Init("info")
	output := captureOutput(func() {
		Infof("user %s logged in", "alice")
	})
	assert.Contains(t, output, "user alice logged in")
}

func TestWarnf(t *testing.T) {
	Init("warn")
	output := captureOutput(func() {
		Warnf("retry %d of %d", 1, 3)
	})
	assert.Contains(t, output, "retry 1 of 3")
}

func TestErrorf(t *testing.T) {
	Init("error")

	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	errorLogger = log.New(w, "[ERROR] ", log.Ldate|log.Ltime)

	Errorf("connection failed: %v", "timeout")

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	assert.Contains(t, output, "connection failed: timeout")
}
