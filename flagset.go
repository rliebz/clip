package clip

import (
	"encoding"
	"fmt"
	"io"
	"iter"
	"os"
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
		env []string,
	)

	// DefineString creates a new string flag.
	DefineString(
		p *string,
		name string,
		short string,
		value string,
		usage string,
		env []string,
	)

	// DefineString creates a new text flag.
	DefineText(
		p encoding.TextUnmarshaler,
		name string,
		short string,
		value encoding.TextMarshaler,
		usage string,
		env []string,
	)

	// Has returns whether a flag exists by name.
	Has(name string) bool

	// HasShort returns whether a flag exists by short name.
	HasShort(name string) bool

	// Parse parses a set of command-line arguments.
	Parse(args []string) error

	// private prevents custom implementations
	private()
}

// NewFlagSet returns a FlagSet.
func NewFlagSet(name string) FlagSet {
	pfs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	pfs.SetOutput(io.Discard)

	return &flagSetImpl{
		flagSet:   pfs,
		nameToEnv: make(map[string][]string),
	}
}

type flagSetImpl struct {
	flagSet   *pflag.FlagSet
	nameToEnv map[string][]string
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
	env []string,
) {
	fs.flagSet.BoolVarP(p, name, short, value, usage)
	fs.nameToEnv[name] = env
}

// DefineString defines a string flag.
func (fs *flagSetImpl) DefineString(
	p *string,
	name string,
	short string,
	value string,
	usage string,
	env []string,
) {
	fs.flagSet.StringVarP(p, name, short, value, usage)
	fs.nameToEnv[name] = env
}

// DefineText defines a flag based on [encoding.TextMarshaler] and
// [encoding.TextUnmarshaler].
func (fs *flagSetImpl) DefineText(
	p encoding.TextUnmarshaler,
	name string,
	short string,
	value encoding.TextMarshaler,
	usage string,
	env []string,
) {
	fs.flagSet.TextVarP(p, name, short, value, usage)
	fs.nameToEnv[name] = env
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

	if err := fs.flagSet.Parse(args); err != nil {
		return err
	}

	for f := range fs.all() {
		if err := fs.parseEnv(f); err != nil {
			return err
		}
	}

	return nil
}

// all returns an iterator over all registered flags.
func (fs *flagSetImpl) all() iter.Seq[*pflag.Flag] {
	return func(yield func(*pflag.Flag) bool) {
		done := false

		fs.flagSet.VisitAll(func(f *pflag.Flag) {
			if done {
				return
			}

			done = !yield(f)
		})
	}
}

func (fs *flagSetImpl) parseEnv(f *pflag.Flag) error {
	if f.Changed {
		return nil
	}

	for _, env := range fs.nameToEnv[f.Name] {
		if v, ok := os.LookupEnv(env); ok {
			return fs.flagSet.Set(f.Name, v)
		}
	}

	return nil
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

func (fs *flagSetImpl) private() {}

func isFlag(arg string) bool {
	return strings.HasPrefix(arg, "-")
}
