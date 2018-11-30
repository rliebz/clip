package clipflag

import "github.com/spf13/pflag"

// Flag is an immutable command-line flag.
type Flag struct {
	name        string
	short       string
	summary     string
	description string
	define      func(*pflag.FlagSet)
	defineShort func(*pflag.FlagSet)

	// TODO: These
	envVar     string // nolint
	deprecated bool   // nolint
	hidden     bool   // nolint
}

// Name returns the name of the flag.
func (f *Flag) Name() string { return f.name }

// Short returns a single character flag name.
func (f *Flag) Short() string { return f.short }

// Summary returns a one-line description of the flag.
func (f *Flag) Summary() string { return f.summary }

// Description returns a multi-line description of the command.
func (f *Flag) Description() string { return f.description }

// Define attaches a flag to a flagset.
func (f *Flag) Define(fs *pflag.FlagSet) {
	if f.short == "" {
		f.define(fs)
	} else {
		f.defineShort(fs)
	}
}
