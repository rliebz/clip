package clip

// Flag is the interface for any flag.
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
type FlagSet interface {
	// Args returns the non-flag arguments passed.
	Args() []string

	// Changed returns whether a flag was explicitly passed to change its value.
	Changed(name string) bool

	// DefineBool creates a new boolean flag.
	DefineBool(p *bool, name string, short string, value bool, usage string)

	// DefineString creates a new string flag.
	DefineString(p *string, name string, short string, value string, usage string)

	// Has returns whether a flag exists by name.
	Has(name string) bool

	// HasShort returns whether a flag exists by short name.
	HasShort(name string) bool

	// Parse parses a set of command-line arguments.
	Parse(args []string) error
}
