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

// TextVar can be marshaled to and from text.
//
// This is the recommended interface to implement for custom types.
type TextVar interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

// FlagSet is the interface for a set of flags.
type FlagSet interface {
	// Args returns the non-flag arguments passed.
	Args() []string

	// Changed returns whether a flag was explicitly passed to change its value.
	Changed(name string) bool

	// DefineBool creates a new boolean flag.
	DefineBool(*bool, *Flag)

	// DefineString creates a new string flag.
	DefineString(*string, *Flag)

	// DefineString creates a new text flag.
	DefineTextVar(TextVar, *Flag)

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
func (fs *flagSetImpl) DefineBool(p *bool, f *Flag) {
	fs.flagSet.BoolVarP(p, f.name, f.short, *p, f.summary)
	fs.nameToEnv[f.name] = f.env
}

// DefineString defines a string flag.
func (fs *flagSetImpl) DefineString(p *string, f *Flag) {
	fs.flagSet.StringVarP(p, f.name, f.short, *p, f.summary)
	fs.nameToEnv[f.name] = f.env
}

// DefineTextVar defines a flag based on [encoding.TextMarshaler] and
// [encoding.TextUnmarshaler].
func (fs *flagSetImpl) DefineTextVar(p TextVar, f *Flag) {
	fs.flagSet.TextVarP(p, f.name, f.short, p, f.summary)
	fs.nameToEnv[f.name] = f.env
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
