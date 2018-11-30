package clipflag

import "github.com/spf13/pflag"

// NewToggle creates a new toggle flag.
// Toggle flags have no associated value, but can be passed like boolean flags
// to toggle something on. This is the simplest way to create an action flag.
func NewToggle(name string, options ...func(*Flag)) *Flag {
	f := newConfig(name, options...)

	f.define = func(fs *pflag.FlagSet) {
		fs.Bool(name, false, f.summary)
	}

	f.defineShort = func(fs *pflag.FlagSet) {
		fs.BoolP(name, f.short, false, f.summary)
	}

	return &f
}

// NewBool creates a new boolean flag.
func NewBool(value *bool, name string, options ...func(*Flag)) *Flag {
	f := newConfig(name, options...)

	f.define = func(fs *pflag.FlagSet) {
		fs.BoolVar(value, name, *value, f.summary)
	}

	f.defineShort = func(fs *pflag.FlagSet) {
		fs.BoolVarP(value, name, f.short, *value, f.summary)
	}

	return &f
}

// NewString creates a new string flag.
func NewString(value *string, name string, options ...func(*Flag)) *Flag {
	f := newConfig(name, options...)

	f.define = func(fs *pflag.FlagSet) {
		fs.StringVar(value, name, *value, f.summary)
	}

	f.defineShort = func(fs *pflag.FlagSet) {
		fs.StringVarP(value, name, f.short, *value, f.summary)
	}

	return &f
}

func newConfig(name string, options ...func(*Flag)) Flag {
	f := Flag{name: name}
	for _, o := range options {
		o(&f)
	}

	return f
}

// WithShort adds a short name to a flag.
// Panics if the name is not exactly one ASCII character.
func WithShort(name string) func(*Flag) {
	return func(f *Flag) {
		f.short = name
	}
}

// WithSummary adds a one-line description to a flag.
func WithSummary(summary string) func(*Flag) {
	return func(f *Flag) {
		f.summary = summary
	}
}

// WithDescription adds a multi-line description to a flag.
func WithDescription(description string) func(*Flag) {
	return func(f *Flag) {
		f.description = description
	}
}
