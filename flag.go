package clip

import (
	"cmp"
	"fmt"
)

// flagDef is an immutable command-line flag.
//
// Methods are defined for use in help text.
type flagDef struct {
	name        string
	short       string
	description string
	deprecated  string
	hidden      bool
	action      func(*Context) error
	env         []string
	placeholder string

	boolVal string
	setFunc func(string) error
	changed bool
}

// Usage returns padded usage text for use in help docs.
func (f *flagDef) Usage() string {
	usage := "      "
	if f.short != "" {
		usage = "  -" + f.short + ", "
	}

	usage += "--" + f.name

	if f.placeholder != "" {
		sep := " "
		if f.boolVal != "" {
			sep = "="
		}
		usage += sep + f.placeholder
	}

	return usage
}

// Description returns a description of the flag.
func (f *flagDef) Description() string { return f.description }

// Deprecated returns the deprecation message, if deprecated.
func (f *flagDef) Deprecated() string { return f.deprecated }

// Env returns the list of environment variables.
func (f *flagDef) Env() []string {
	return f.env
}

// Hidden returns whether a flag should be hidden from help and tab completion.
func (f *flagDef) Hidden() bool { return f.hidden }

// set assigns a string value to a flag.
func (f *flagDef) set(v string) error {
	if err := f.setFunc(v); err != nil {
		return err
	}
	f.changed = true
	return nil
}

// ToggleFlag creates a new toggle flag.
//
// Toggle flags have no associated value, but can be passed like boolean flags
// to toggle something on. This is the simplest way to create an action flag.
func ToggleFlag(name string, options ...FlagOption) CommandOption {
	return func(c *commandConfig) {
		f := newFlag(name, options...)
		f.boolVal = "true"
		f.setFunc = func(s string) error {
			switch s {
			case "true", "1":
				return nil
			default:
				return fmt.Errorf("invalid toggle flag option: %s", s)
			}
		}

		c.addFlag(f)
	}
}

// BoolFlag creates a new boolean flag.
func BoolFlag(value *bool, name string, options ...FlagOption) CommandOption {
	return func(c *commandConfig) {
		f := newFlag(name, options...)
		f.boolVal = "true"
		f.setFunc = func(s string) error {
			switch s {
			case "true", "1":
				*value = true
			case "false", "0":
				*value = false
			default:
				return fmt.Errorf("non-boolean value: %s", s)
			}

			return nil
		}

		c.addFlag(f)
	}
}

// StringFlag creates a new string flag.
func StringFlag(value *string, name string, options ...FlagOption) CommandOption {
	return func(c *commandConfig) {
		f := newFlag(name, options...)
		f.placeholder = cmp.Or(f.placeholder, "<string>")
		f.setFunc = func(s string) error {
			*value = s
			return nil
		}

		c.addFlag(f)
	}
}

// TextVarFlag creates a new flag based on [encoding.TextMarshaler] and
// [encoding.TextUnmarshaler].
func TextVarFlag(value TextVar, name string, options ...FlagOption) CommandOption {
	return func(c *commandConfig) {
		f := newFlag(name, options...)
		f.placeholder = cmp.Or(f.placeholder, "<value>")
		f.setFunc = func(s string) error {
			return value.UnmarshalText([]byte(s))
		}

		c.addFlag(f)
	}
}

// FlagOption is an option for creating a Flag.
type FlagOption func(*flagConfig)

type flagConfig struct {
	short       string
	description string
	deprecated  string
	env         []string
	placeholder string
	hidden      bool
	action      func(*Context) error
}

func newFlag(name string, options ...FlagOption) *flagDef {
	c := flagConfig{}
	for _, o := range options {
		o(&c)
	}

	return &flagDef{
		name:        name,
		short:       c.short,
		description: c.description,
		env:         c.env,
		deprecated:  c.deprecated,
		placeholder: c.placeholder,
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

// FlagDescription adds a description to a flag.
//
// Descriptions can span multiple lines.
func FlagDescription(description string) FlagOption {
	return func(c *flagConfig) {
		c.description = description
	}
}

// FlagDeprecated adds a deprecation to a flag.
func FlagDeprecated(derepcation string) FlagOption {
	return func(c *flagConfig) {
		c.deprecated = derepcation
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

// FlagPlaceholder sets the name of the flag's placeholder value.
func FlagPlaceholder(name string) FlagOption {
	return func(c *flagConfig) {
		c.placeholder = "<" + name + ">"
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
