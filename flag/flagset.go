package flag

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/pflag"

	"github.com/rliebz/clip"
)

// NewFlagSet returns a FlagSet.
func NewFlagSet(name string) *FlagSet {
	pfs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	pfs.SetOutput(io.Discard)

	return &FlagSet{
		flagSet: pfs,
	}
}

// FlagSet represents a set of defined flags.
type FlagSet struct { //nolint:revive
	flagSet *pflag.FlagSet
}

var _ clip.FlagSet = (*FlagSet)(nil)

// Args returns non-flag arguments.
func (fs *FlagSet) Args() []string {
	return fs.flagSet.Args()
}

// Changed returns whether a variable was changed.
func (fs *FlagSet) Changed(name string) bool {
	return fs.flagSet.Changed(name)
}

// DefineBool defines a bool flag.
func (fs *FlagSet) DefineBool(p *bool, name string, short string, value bool, usage string) {
	fs.flagSet.BoolVarP(p, name, short, value, usage)
}

// DefineString defines a string flag.
func (fs *FlagSet) DefineString(p *string, name string, short string, value string, usage string) {
	fs.flagSet.StringVarP(p, name, short, value, usage)
}

// Has returns whether a flagset has a flag by a name.
func (fs *FlagSet) Has(name string) bool {
	return fs.flagSet.Lookup(name) != nil
}

// HasShort returns whether a flagset has a flag by a short name.
func (fs *FlagSet) HasShort(name string) bool {
	return fs.flagSet.ShorthandLookup(name) != nil
}

// Parse a set of command-line arguments as flags.
func (fs *FlagSet) Parse(args []string) error {
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
func (fs *FlagSet) nextArgIndex(args []string) (int, error) {
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

func (fs *FlagSet) getFlagFromArg(arg string) (*pflag.Flag, error) {
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
