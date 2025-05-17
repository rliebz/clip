package clip

import (
	"encoding"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/pflag"
)

// FlagSet is the interface for a set of flags.
type FlagSet interface {
	// Args returns the non-flag arguments passed.
	Args() []string

	// Changed returns whether a flag was explicitly passed to change its value.
	Changed(name string) bool

	// DefineBool creates a new boolean flag.
	DefineBool(
		p *bool,
		name string,
		short string,
		value bool,
		usage string,
	)

	// DefineString creates a new string flag.
	DefineString(
		p *string,
		name string,
		short string,
		value string,
		usage string,
	)

	// DefineString creates a new text flag.
	DefineText(
		p encoding.TextUnmarshaler,
		name string,
		short string,
		value encoding.TextMarshaler,
		usage string,
	)

	// Has returns whether a flag exists by name.
	Has(name string) bool

	// HasShort returns whether a flag exists by short name.
	HasShort(name string) bool

	// Parse parses a set of command-line arguments.
	Parse(args []string) error
}

// NewFlagSet returns a FlagSet.
func NewFlagSet(name string) FlagSet {
	pfs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	pfs.SetOutput(io.Discard)

	return &flagSetImpl{
		flagSet: pfs,
	}
}

type flagSetImpl struct {
	flagSet *pflag.FlagSet
}

// Args returns non-flag arguments.
func (fs *flagSetImpl) Args() []string {
	return fs.flagSet.Args()
}

// Changed returns whether a variable was changed.
func (fs *flagSetImpl) Changed(name string) bool {
	return fs.flagSet.Changed(name)
}

// DefineBool defines a bool flag.
func (fs *flagSetImpl) DefineBool(
	p *bool,
	name string,
	short string,
	value bool,
	usage string,
) {
	fs.flagSet.BoolVarP(p, name, short, value, usage)
}

// DefineString defines a string flag.
func (fs *flagSetImpl) DefineString(
	p *string,
	name string,
	short string,
	value string,
	usage string,
) {
	fs.flagSet.StringVarP(p, name, short, value, usage)
}

// DefineText defines a flag based on [encoding.TextMarshaler]/[encoding.TextUnmarshaler].
func (fs *flagSetImpl) DefineText(
	p encoding.TextUnmarshaler,
	name string,
	short string,
	value encoding.TextMarshaler,
	usage string,
) {
	fs.flagSet.TextVarP(p, name, short, value, usage)
}

// Has returns whether a flagset has a flag by a name.
func (fs *flagSetImpl) Has(name string) bool {
	return fs.flagSet.Lookup(name) != nil
}

// HasShort returns whether a flagset has a flag by a short name.
func (fs *flagSetImpl) HasShort(name string) bool {
	return fs.flagSet.ShorthandLookup(name) != nil
}

// Parse a set of command-line arguments as flags.
func (fs *flagSetImpl) Parse(args []string) error {
	i, err := fs.nextArgIndex(args)
	if err != nil {
		return err
	}

	if i != -1 {
		args = append(args[:i], append([]string{"--"}, args[i:]...)...)
	}

	return fs.flagSet.Parse(args)
}

// nextArgIndex finds the index of the next arg in the arg list.
// If no args are present, -1 is returned.
func (fs *flagSetImpl) nextArgIndex(args []string) (int, error) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !isFlag(arg) || arg == "--" {
			return i, nil
		}

		f, err := fs.getFlagFromArg(arg)
		if err != nil {
			return 0, err
		}

		if !strings.Contains(arg, "=") && f.Value.Type() != "bool" {
			i++
		}
	}

	return -1, nil
}

func (fs *flagSetImpl) getFlagFromArg(arg string) (*pflag.Flag, error) {
	fname := strings.SplitN(arg, "=", 2)[0]

	if strings.HasPrefix(fname, "--") {
		fname = strings.TrimPrefix(fname, "--")
		if f := fs.flagSet.Lookup(fname); f != nil {
			return f, nil
		}
		return nil, fmt.Errorf("unknown flag: %s", fname)
	}

	fname = fname[len(fname)-1:]
	if f := fs.flagSet.ShorthandLookup(fname); f != nil {
		return f, nil
	}
	return nil, fmt.Errorf("unknown shorthand flag: '%s' in %s", fname, arg)
}

func isFlag(arg string) bool {
	return strings.HasPrefix(arg, "-")
}
