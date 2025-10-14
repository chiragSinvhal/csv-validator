package logger

import (
	"log"
	"os"
	"strings"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var (
	currentLevel LogLevel = INFO
	debugLogger           = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger            = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger            = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger           = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
)

// Init initializes the logger with the specified level
func Init(level string) {
	switch strings.ToLower(level) {
	case "debug":
		currentLevel = DEBUG
	case "info":
		currentLevel = INFO
	case "warn", "warning":
		currentLevel = WARN
	case "error":
		currentLevel = ERROR
	default:
		currentLevel = INFO
	}
}

// Debug logs a debug message
func Debug(message string) {
	if currentLevel <= DEBUG {
		debugLogger.Println(message)
	}
}

// Info logs an info message
func Info(message string) {
	if currentLevel <= INFO {
		infoLogger.Println(message)
	}
}

// Warn logs a warning message
func Warn(message string) {
	if currentLevel <= WARN {
		warnLogger.Println(message)
	}
}

// Error logs an error message
func Error(message string) {
	if currentLevel <= ERROR {
		errorLogger.Println(message)
	}
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	if currentLevel <= DEBUG {
		debugLogger.Printf(format, args...)
	}
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	if currentLevel <= INFO {
		infoLogger.Printf(format, args...)
	}
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	if currentLevel <= WARN {
		warnLogger.Printf(format, args...)
	}
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	if currentLevel <= ERROR {
		errorLogger.Printf(format, args...)
	}
}
