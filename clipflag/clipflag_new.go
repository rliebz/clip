package clipflag

import "github.com/spf13/pflag"

// NewBool creates a new boolean flag.
func NewBool(value *bool, name string, options ...func(*flag)) *flag {
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
func NewString(value *string, name string, options ...func(*flag)) *flag {
	f := newConfig(name, options...)

	f.define = func(fs *pflag.FlagSet) {
		fs.StringVar(value, name, *value, f.summary)
	}

	f.defineShort = func(fs *pflag.FlagSet) {
		fs.StringVarP(value, name, f.short, *value, f.summary)
	}

	return &f
}

func newConfig(name string, options ...func(*flag)) flag {
	f := flag{name: name}
	for _, o := range options {
		o(&f)
	}

	return f
}

// WithShort adds a short name to a flag.
// Panics if the name is not exactly one ASCII character.
func WithShort(name string) func(*flag) {
	return func(f *flag) {
		f.short = name
	}
}

// WithSummary adds a one-line description to a flag.
func WithSummary(summary string) func(*flag) {
	return func(f *flag) {
		f.summary = summary
	}
}

// WithDescription adds a multi-line description to a flag.
func WithDescription(description string) func(*flag) {
	return func(f *flag) {
		f.description = description
	}
}
