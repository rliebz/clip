package clip

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// Flag is the interface for any flag.
// Typically, this will be implemented by a flag from the clipflag package.
type Flag interface {
	Name() string
	Summary() string

	// TODO: Replace pflag.FlagSet with an interface
	Define(*pflag.FlagSet)
}

// WithFlag adds a flag.
func WithFlag(f Flag) func(*Command) {
	return func(cmd *Command) {
		f.Define(cmd.flagSet)
	}
}

// WithActionFlag adds a flag that performs an action and nothing else.
// Flags such as `--help` or `--version` fall under this category.
func WithActionFlag(f Flag, action func(*Context) error) func(*Command) {
	return func(cmd *Command) {
		oldAction := cmd.flagAction
		f.Define(cmd.flagSet)
		cmd.flagAction = func(ctx *Context) (bool, error) {
			if wasSet, err := oldAction(ctx); wasSet {
				return true, err
			}
			if cmd.flagSet.Changed(f.Name()) {
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
