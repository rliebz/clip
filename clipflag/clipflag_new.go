package clipflag

import "github.com/spf13/pflag"

// NewToggle creates a new toggle flag.
// Toggle flags have no associated value, but can be passed like boolean flags
// to toggle something on. This is the simplest way to create an action flag.
func NewToggle(name string, options ...Option) *Flag {
	f := newFlag(name, options...)

	f.define = func(fs *pflag.FlagSet) {
		fs.Bool(name, false, f.summary)
	}

	f.defineShort = func(fs *pflag.FlagSet) {
		fs.BoolP(name, f.short, false, f.summary)
	}

	return &f
}

// NewBool creates a new boolean flag.
func NewBool(value *bool, name string, options ...Option) *Flag {
	f := newFlag(name, options...)

	f.define = func(fs *pflag.FlagSet) {
		fs.BoolVar(value, name, *value, f.summary)
	}

	f.defineShort = func(fs *pflag.FlagSet) {
		fs.BoolVarP(value, name, f.short, *value, f.summary)
	}

	return &f
}

// NewString creates a new string flag.
func NewString(value *string, name string, options ...Option) *Flag {
	f := newFlag(name, options...)

	f.define = func(fs *pflag.FlagSet) {
		fs.StringVar(value, name, *value, f.summary)
	}

	f.defineShort = func(fs *pflag.FlagSet) {
		fs.StringVarP(value, name, f.short, *value, f.summary)
	}

	return &f
}

// Option is an option for creating a Flag.
type Option func(*config)

type config struct {
	short       string
	summary     string
	description string
	envVar      string
	deprecated  bool
	hidden      bool
}

func newFlag(name string, options ...Option) Flag {
	c := config{}
	for _, o := range options {
		o(&c)
	}

	return Flag{
		name:        name,
		short:       c.short,
		summary:     c.summary,
		description: c.description,
		envVar:      c.envVar,
		deprecated:  c.deprecated,
		hidden:      c.hidden,
	}
}

// AsHidden prevents the flag from being shown.
func AsHidden(c *config) {
	c.hidden = true
}

// WithShort adds a short name to a flag.
// Panics if the name is not exactly one ASCII character.
func WithShort(name string) Option {
	return func(c *config) {
		c.short = name
	}
}

// WithSummary adds a one-line description to a flag.
func WithSummary(summary string) Option {
	return func(c *config) {
		c.summary = summary
	}
}

// WithDescription adds a multi-line description to a flag.
func WithDescription(description string) Option {
	return func(c *config) {
		c.description = description
	}
}
