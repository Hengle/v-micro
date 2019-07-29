// Package log is a global internal logger
package log

// Logger logger interface
type Logger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

var (
	// the local logger
	logger Logger
)

// Log log
func Log(v ...interface{}) {
	logger.Log(v...)
}

// Logf logf
func Logf(format string, v ...interface{}) {
	logger.Logf(format, v...)
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
