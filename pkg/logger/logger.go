package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"time"
)

var log *logrus.Logger

// LogConfig holds the configuration for the logger
type LogConfig struct {
	Level      string
	Format     string // "json" or "text"
	OutputPath string
}

// InitLogger initializes the logger with the given configuration
func InitLogger(config LogConfig) error {
	log = logrus.New()
	
	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return err
	}
	log.SetLevel(level)
	
	// Set output format
	if config.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}
	
	// Set output
	if config.OutputPath != "" {
		file, err := os.OpenFile(config.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		log.SetOutput(file)
	} else {
		log.SetOutput(os.Stdout)
	}
	
	return nil
}

// WithContext adds common fields to the log entry
func WithContext(fields logrus.Fields) *logrus.Entry {
	// Add common fields like goroutine ID, timestamp, etc.
	fields["goroutine_id"] = runtime.NumGoroutine()
	return log.WithFields(fields)
}

// Info logs an info level message with structured fields
func Info(message string, fields logrus.Fields) {
	WithContext(fields).Info(message)
}

// Warning logs a warning level message with structured fields
func Warning(message string, fields logrus.Fields) {
	WithContext(fields).Warn(message)
}

// Error logs an error level message with structured fields
func Error(message string, fields logrus.Fields) {
	WithContext(fields).Error(message)
}

// Fatal logs a fatal level message with structured fields and exits
func Fatal(message string, fields logrus.Fields) {
	WithContext(fields).Fatal(message)
}

// Debug logs a debug level message with structured fields
func Debug(message string, fields logrus.Fields) {
	WithContext(fields).Debug(message)
}
