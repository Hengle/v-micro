package log

// DefaultLogger default logger
var DefaultLogger Logger

// Info info
func Info(v ...interface{}) {
	DefaultLogger.Info(v...)
}

// Infof infof
func Infof(format string, v ...interface{}) {
	DefaultLogger.Infof(format, v...)
}

// Error error
func Error(v ...interface{}) {
	DefaultLogger.Error(v...)
}

// Errorf errorf
func Errorf(format string, v ...interface{}) {
	DefaultLogger.Errorf(format, v...)
}

// Fatal fatal
func Fatal(v ...interface{}) {
	DefaultLogger.Fatal(v...)
}

// Fatalf fatalf
func Fatalf(format string, v ...interface{}) {
	DefaultLogger.Fatalf(format, v...)
}

// SetLogger sets the local logger
func SetLogger(l Logger) {
	DefaultLogger = l
}

// GetLogger returns the local logger
func GetLogger() Logger {
	return DefaultLogger
}
