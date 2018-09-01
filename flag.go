package clip

import (
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