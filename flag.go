package clip

import "encoding"

// Flag is the interface for any flag.
type Flag interface {
	Name() string
	Short() string
	Summary() string
	Hidden() bool

	// Attach the flag to a given flagset.
	//
	// This method is invoked before flags are parsed.
	Attach(FlagSet)
}

// flagImpl is an immutable command-line flag.
type flagImpl struct {
	name        string
	short       string
	summary     string
	description string
	hidden      bool
	attach      func(FlagSet)

	// TODO: Help text
	env []string

	// TODO: This
	deprecated bool
}

// Name returns the name of the flag.
func (f *flagImpl) Name() string { return f.name }

// Short returns a single character flag name.
func (f *flagImpl) Short() string { return f.short }

// Summary returns a one-line description of the flag.
func (f *flagImpl) Summary() string { return f.summary }

// Description returns a multi-line description of the flag.
func (f *flagImpl) Description() string { return f.description }

// Hidden returns whether a flag should be hidden from help and tab completion.
func (f *flagImpl) Hidden() bool { return f.hidden }

// Attach attaches a flag to a flagset.
func (f *flagImpl) Attach(fs FlagSet) {
	f.attach(fs)
}

// NewToggleFlag creates a new toggle flag.
// Toggle flags have no associated value, but can be passed like boolean flags
// to toggle something on. This is the simplest way to create an action flag.
func NewToggleFlag(name string, options ...FlagOption) Flag {
	f := newFlag(name, options...)

	f.attach = func(fs FlagSet) {
		p := new(bool)
		fs.DefineBool(p, name, f.short, false, f.summary, f.env)
	}

	return &f
}

// NewBoolFlag creates a new boolean flag.
func NewBoolFlag(value *bool, name string, options ...FlagOption) Flag {
	f := newFlag(name, options...)

	f.attach = func(fs FlagSet) {
		fs.DefineBool(value, name, f.short, *value, f.summary, f.env)
	}

	return &f
}

// NewStringFlag creates a new string flag.
func NewStringFlag(value *string, name string, options ...FlagOption) Flag {
	f := newFlag(name, options...)

	f.attach = func(fs FlagSet) {
		fs.DefineString(value, name, f.short, *value, f.summary, f.env)
	}

	return &f
}

// NewTextVarFlag creates a new flag based on [encoding.TextMarshaler] and
// [encoding.TextUnmarshaler].
func NewTextVarFlag(
	value interface {
		encoding.TextMarshaler
		encoding.TextUnmarshaler
	},
	name string,
	options ...FlagOption,
) Flag {
	f := newFlag(name, options...)

	f.attach = func(fs FlagSet) {
		fs.DefineTextVar(value, name, f.short, value, f.summary, f.env)
	}

	return &f
}

// FlagOption is an option for creating a Flag.
type FlagOption func(*flagConfig)

type flagConfig struct {
	short       string
	summary     string
	description string
	env         []string
	deprecated  bool
	hidden      bool
}

func newFlag(name string, options ...FlagOption) flagImpl {
	c := flagConfig{}
	for _, o := range options {
		o(&c)
	}

	return flagImpl{
		name:        name,
		short:       c.short,
		summary:     c.summary,
		description: c.description,
		env:         c.env,
		deprecated:  c.deprecated,
		hidden:      c.hidden,
	}
}

// FlagHidden prevents the flag from being shown.
func FlagHidden(c *flagConfig) {
	c.hidden = true
}

// FlagShort adds a short name to a flag.
// Panics if the name is not exactly one ASCII character.
func FlagShort(name string) FlagOption {
	return func(c *flagConfig) {
		c.short = name
	}
}

// FlagSummary adds a one-line description to a flag.
func FlagSummary(summary string) FlagOption {
	return func(c *flagConfig) {
		c.summary = summary
	}
}

// FlagDescription adds a multi-line description to a flag.
func FlagDescription(description string) FlagOption {
	return func(c *flagConfig) {
		c.description = description
	}
}

// FlagEnv sets the list of environment variables for a flag.
//
// Successive calls will replace earlier values.
func FlagEnv(env ...string) FlagOption {
	return func(c *flagConfig) {
		c.env = env
	}
}
