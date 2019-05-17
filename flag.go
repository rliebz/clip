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
// Typically, this will be implemented by clipflag.
type FlagSet interface {
	Args() []string
	Changed(name string) bool
	DefineBool(p *bool, name string, short string, value bool, usage string)
	DefineString(p *string, name string, short string, value string, usage string)
	Has(name string) bool
	HasShort(name string) bool
	Parse(args []string) error
}
