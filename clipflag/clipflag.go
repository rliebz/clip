package clipflag

import (
	"github.com/rliebz/clip"
)

// Flag is an immutable command-line flag.
type Flag struct {
	name        string
	short       string
	summary     string
	description string
	hidden      bool
	define      func(clip.FlagSet)

	// TODO: These
	envVar     string // nolint
	deprecated bool   // nolint
}

var _ clip.Flag = (*Flag)(nil)

// Name returns the name of the flag.
func (f *Flag) Name() string { return f.name }

// Short returns a single character flag name.
func (f *Flag) Short() string { return f.short }

// Summary returns a one-line description of the flag.
func (f *Flag) Summary() string { return f.summary }

// Description returns a multi-line description of the command.
func (f *Flag) Description() string { return f.description }

// Hidden returns whether a flag should be hidden from help and tab completion.
func (f *Flag) Hidden() bool { return f.hidden }

// Define attaches a flag to a flagset.
func (f *Flag) Define(fs clip.FlagSet) {
	f.define(fs)
}
