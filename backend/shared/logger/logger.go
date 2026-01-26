package logger

import "log"

// Logger provides structured logging capabilities (Datadog-ready)
type Logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, err error, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Debug(msg string, fields map[string]interface{})
}

// DefaultLogger is a simple implementation
type DefaultLogger struct{}

func (l *DefaultLogger) Info(msg string, fields map[string]interface{}) {
	log.Printf("[INFO] %s %+v", msg, fields)
}

func (l *DefaultLogger) Error(msg string, err error, fields map[string]interface{}) {
	log.Printf("[ERROR] %s: %v %+v", msg, err, fields)
}

func (l *DefaultLogger) Warn(msg string, fields map[string]interface{}) {
	log.Printf("[WARN] %s %+v", msg, fields)
}

func (l *DefaultLogger) Debug(msg string, fields map[string]interface{}) {
	log.Printf("[DEBUG] %s %+v", msg, fields)
}
