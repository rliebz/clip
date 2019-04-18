package command

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// Flag is the interface for any flag.
// Typically, this will be implemented by a flag from the clipflag package.
type Flag interface {
	Name() string
	Short() string
	Summary() string
	Hidden() bool

	// Define adds the flag to a given flagset.
	// This method is invoked when creating a new command before the flags are
	// parsed.
	//
	// TODO: Replace pflag.FlagSet with an interface
	Define(*pflag.FlagSet)
}

// WithFlag adds a flag.
// Typically, flags from the clipflag package will be passed here.
func WithFlag(f Flag) Option {
	return func(c *config) {
		f.Define(c.flagSet)
		if !f.Hidden() {
			c.visibleFlags = append(c.visibleFlags, f)
		}
	}
}

// WithActionFlag adds a flag that performs an action and nothing else.
// Flags such as --help or --version fall under this category.
//
// The action will occur if the flag is passed, regardless of the value, so
// typically clipflag.NewToggle will be used here.
func WithActionFlag(f Flag, action func(*Context) error) Option {
	return func(c *config) {
		oldAction := c.flagAction
		f.Define(c.flagSet)
		if !f.Hidden() {
			c.visibleFlags = append(c.visibleFlags, f)
		}
		c.flagAction = func(ctx *Context) (bool, error) {
			if wasSet, err := oldAction(ctx); wasSet {
				return true, err
			}
			if c.flagSet.Changed(f.Name()) {
				return true, action(ctx)
			}
			return false, nil
		}
	}
}

func parse(ctx *Context, args []string) error {
	i, err := splitAtFirstArg(ctx, args)
	if err != nil {
		return err
	}

	args = append(args[:i], append([]string{"--"}, args[i:]...)...)
	return ctx.flagSet.Parse(args)
}

func splitAtFirstArg(ctx *Context, args []string) (int, error) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !isFlag(arg) || arg == "--" {
			return i, nil
		}

		f, err := getFlagFromArg(ctx, arg)
		if err != nil {
			return 0, err
		}

		if !strings.Contains(arg, "=") && f.Value.Type() != "bool" {
			i++
		}
	}

	return len(args), nil
}

func isFlag(arg string) bool {
	return strings.HasPrefix(arg, "-")
}

func getFlagFromArg(ctx *Context, arg string) (*pflag.Flag, error) {
	fname := strings.SplitN(arg, "=", 2)[0]

	if strings.HasPrefix(fname, "--") {
		fname = strings.TrimPrefix(fname, "--")
		if f := ctx.flagSet.Lookup(fname); f != nil {
			return f, nil
		}
		return nil, fmt.Errorf("unknown flag: %s", fname)
	}

	fname = fname[len(fname)-1:]
	if f := ctx.flagSet.ShorthandLookup(fname); f != nil {
		return f, nil
	}
	return nil, fmt.Errorf("unknown shorthand flag: '%s' in %s", fname, arg)
}
