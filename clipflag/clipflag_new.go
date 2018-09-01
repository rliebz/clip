package clipflag

import "github.com/spf13/pflag"

// NewBool creates a new boolean flag.
func NewBool(value *bool, name string, options ...func(*flag)) *flag {
	f := newConfig(name, options...)

	f.define = func(fs *pflag.FlagSet) {
		fs.BoolVar(value, name, *value, f.summary)
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
