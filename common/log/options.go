package log

import "context"

// Options options
type Options struct {
	Name  string
	Level LevelType

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Option option
type Option func(*Options)

// Name log name
func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// Level level
func Level(level LevelType) Option {
	return func(o *Options) {
		o.Level = level
	}
}
