// Package log is a global internal logger
package log

// LevelType level type
type LevelType int

const (
	// TraceLevel trace level
	TraceLevel LevelType = iota
	// DebugLevel debug level
	DebugLevel
	// InfoLevel info level
	InfoLevel
	// WarnLevel warn level
	WarnLevel
	// ErrorLevel error level
	ErrorLevel
	// FatalLevel fatal level
	FatalLevel
	// PanicLevel panic level
	PanicLevel
)

// Logger logger interface
type Logger interface {
	Init(opt ...Option) error
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	String() string
}
