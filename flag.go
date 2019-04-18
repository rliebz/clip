package clip

import "github.com/spf13/pflag"

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
