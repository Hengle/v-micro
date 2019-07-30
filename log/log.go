// Package log is a global internal logger
package log

// Logger logger interface
type Logger interface {
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

var (
	// the local logger
	logger Logger
)

// Info info
func Info(v ...interface{}) {
	logger.Info(v...)
}

// Infof infof
func Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

// Error error
func Error(v ...interface{}) {
	logger.Error(v...)
}

// Errorf errorf
func Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

// SetLogger sets the local logger
func SetLogger(l Logger) {
	logger = l
}

// GetLogger returns the local logger
func GetLogger() Logger {
	return logger
}
