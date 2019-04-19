package clip

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
	Define(FlagSet)
}

// FlagSet is the interface for a set of flags.
// Typically, this will be implemented by github.com/spf13/pflag.
type FlagSet interface {
	BoolVarP(p *bool, name string, short string, value bool, usage string)
	StringVarP(p *string, name string, short string, value string, usage string)
}
