package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/galihrivanto/kotak/config"
)

// Level represents logging level
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

var levelNames = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
}

var levelMap = map[string]Level{
	"debug": DebugLevel,
	"info":  InfoLevel,
	"warn":  WarnLevel,
	"error": ErrorLevel,
}

type Format string

const (
	TextFormat Format = "text"
	JSONFormat Format = "json"
)

type Logger struct {
	level  Level
	format Format
}

var defaultLogger = &Logger{
	level:  InfoLevel,
	format: TextFormat,
}

// Configure sets up the default logger
func Configure(cfg config.Logger) {
	level, ok := levelMap[strings.ToLower(cfg.Level)]
	if !ok {
		level = InfoLevel
	}
	defaultLogger.level = level

	format := Format(strings.ToLower(cfg.Format))
	if format != JSONFormat {
		format = TextFormat
	}
	defaultLogger.format = format
}

func (l *Logger) log(level Level, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format(time.RFC3339)
	levelName := levelNames[level]

	if l.format == JSONFormat {
		fmt.Fprintf(os.Stderr, `{"time":"%s","level":"%s","msg":"%s"}`, timestamp, levelName, fmt.Sprintf(msg, args...))
	} else {
		fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", timestamp, levelName, fmt.Sprintf(msg, args...))
	}
}

// Debug logs a debug message
func Debug(msg string, args ...interface{}) {
	defaultLogger.log(DebugLevel, msg, args...)
}

// Info logs an info message
func Info(msg string, args ...interface{}) {
	defaultLogger.log(InfoLevel, msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...interface{}) {
	defaultLogger.log(WarnLevel, msg, args...)
}

// Error logs an error message
func Error(msg string, args ...interface{}) {
	defaultLogger.log(ErrorLevel, msg, args...)
}

// Fatal logs a fatal message and exits the program
func Fatal(msg string, args ...interface{}) {
	defaultLogger.log(ErrorLevel, msg, args...)
	os.Exit(1)
}
