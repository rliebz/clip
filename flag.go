package clip

import "encoding"

// Flag is the interface for any flag.
type Flag interface {
	Name() string
	Short() string
	Summary() string
	Hidden() bool

	// Define adds the flag to a given flagset.
	// This method is invoked when creating a new command before the flags are
	// parsed.
	Define(FlagSet)
}

// flagImpl is an immutable command-line flag.
type flagImpl struct {
	name        string
	short       string
	summary     string
	description string
	hidden      bool
	define      func(FlagSet)

	// TODO: These
	envVar     string
	deprecated bool
}

// Name returns the name of the flag.
func (f *flagImpl) Name() string { return f.name }

// Short returns a single character flag name.
func (f *flagImpl) Short() string { return f.short }

// Summary returns a one-line description of the flag.
func (f *flagImpl) Summary() string { return f.summary }

// Description returns a multi-line description of the command.
func (f *flagImpl) Description() string { return f.description }

// Hidden returns whether a flag should be hidden from help and tab completion.
func (f *flagImpl) Hidden() bool { return f.hidden }

// Define attaches a flag to a flagset.
func (f *flagImpl) Define(fs FlagSet) {
	f.define(fs)
}

// NewToggle creates a new toggle flag.
// Toggle flags have no associated value, but can be passed like boolean flags
// to toggle something on. This is the simplest way to create an action flag.
func NewToggle(name string, options ...FlagOption) Flag {
	f := newFlag(name, options...)

	f.define = func(fs FlagSet) {
		p := new(bool)
		fs.DefineBool(p, name, f.short, false, f.summary)
	}

	return &f
}

// NewBool creates a new boolean flag.
func NewBool(value *bool, name string, options ...FlagOption) Flag {
	f := newFlag(name, options...)

	f.define = func(fs FlagSet) {
		fs.DefineBool(value, name, f.short, *value, f.summary)
	}

	return &f
}

// NewString creates a new string flag.
func NewString(value *string, name string, options ...FlagOption) Flag {
	f := newFlag(name, options...)

	f.define = func(fs FlagSet) {
		fs.DefineString(value, name, f.short, *value, f.summary)
	}

	return &f
}

// NewText creates a new flag based on [encoding.TextMarshaler]/[encoding.TextUnmarshaler].
func NewText(
	value interface {
		encoding.TextMarshaler
		encoding.TextUnmarshaler
	},
	name string,
	options ...FlagOption,
) Flag {
	f := newFlag(name, options...)

	f.define = func(fs FlagSet) {
		fs.DefineText(value, name, f.short, value, f.summary)
	}

	return &f
}

// FlagOption is an option for creating a Flag.
type FlagOption func(*config)

type config struct {
	short       string
	summary     string
	description string
	envVar      string
	deprecated  bool
	hidden      bool
}

func newFlag(name string, options ...FlagOption) flagImpl {
	c := config{}
	for _, o := range options {
		o(&c)
	}

	return flagImpl{
		name:        name,
		short:       c.short,
		summary:     c.summary,
		description: c.description,
		envVar:      c.envVar,
		deprecated:  c.deprecated,
		hidden:      c.hidden,
	}
}

// FlagHidden prevents the flag from being shown.
func FlagHidden(c *config) {
	c.hidden = true
}

// FlagShort adds a short name to a flag.
// Panics if the name is not exactly one ASCII character.
func FlagShort(name string) FlagOption {
	return func(c *config) {
		c.short = name
	}
}

// FlagSummary adds a one-line description to a flag.
func FlagSummary(summary string) FlagOption {
	return func(c *config) {
		c.summary = summary
	}
}

// WithDescription adds a multi-line description to a flag.
func WithDescription(description string) FlagOption {
	return func(c *config) {
		c.description = description
	}
}
