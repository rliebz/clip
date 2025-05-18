package clip

// Flag is an immutable command-line flag.
type Flag struct {
	name        string
	short       string
	summary     string
	description string
	hidden      bool
	action      func(*Context) error

	// TODO: Help text
	env []string

	// TODO: This
	deprecated bool
}

// Name returns the name of the flag.
func (f *Flag) Name() string { return f.name }

// Short returns a single character flag name.
func (f *Flag) Short() string { return f.short }

// Summary returns a one-line description of the flag.
func (f *Flag) Summary() string { return f.summary }

// Description returns a multi-line description of the flag.
func (f *Flag) Description() string { return f.description }

// Hidden returns whether a flag should be hidden from help and tab completion.
func (f *Flag) Hidden() bool { return f.hidden }

// ToggleFlag creates a new toggle flag.
//
// Toggle flags have no associated value, but can be passed like boolean flags
// to toggle something on. This is the simplest way to create an action flag.
func ToggleFlag(name string, options ...FlagOption) CommandOption {
	return func(c *commandConfig) {
		f := newFlag(name, options...)
		p := new(bool)
		c.flagSet.DefineBool(p, f)
		c.addFlag(f)
	}
}

// BoolFlag creates a new boolean flag.
func BoolFlag(value *bool, name string, options ...FlagOption) CommandOption {
	return func(c *commandConfig) {
		f := newFlag(name, options...)
		c.flagSet.DefineBool(value, f)
		c.addFlag(f)
	}
}

// StringFlag creates a new string flag.
func StringFlag(value *string, name string, options ...FlagOption) CommandOption {
	return func(c *commandConfig) {
		f := newFlag(name, options...)
		c.flagSet.DefineString(value, f)
		c.addFlag(f)
	}
}

// TextVarFlag creates a new flag based on [encoding.TextMarshaler] and
// [encoding.TextUnmarshaler].
func TextVarFlag(value TextVar, name string, options ...FlagOption) CommandOption {
	return func(c *commandConfig) {
		f := newFlag(name, options...)
		c.flagSet.DefineTextVar(value, f)
		c.addFlag(f)
	}
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
	action      func(*Context) error
}

func newFlag(name string, options ...FlagOption) *Flag {
	c := flagConfig{}
	for _, o := range options {
		o(&c)
	}

	return &Flag{
		name:        name,
		short:       c.short,
		summary:     c.summary,
		description: c.description,
		env:         c.env,
		deprecated:  c.deprecated,
		hidden:      c.hidden,
		action:      c.action,
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

// FlagAction sets the behavior of the flag to replace the command's action
// when set.
//
// The action will occur if the flag is passed, regardless of the value.
func FlagAction(action func(*Context) error) FlagOption {
	return func(c *flagConfig) {
		c.action = action
	}
}
