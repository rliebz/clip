package clip

import (
	"encoding"
	"fmt"
	"os"
	"slices"
	"strings"
)

// TextVar can be marshaled to and from text.
//
// This is the recommended interface to implement for custom types.
type TextVar interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

// newFlagSet returns a FlagSet.
func newFlagSet() *flagSet {
	return &flagSet{
		byName:      make(map[string]*flagDef),
		byShortName: make(map[string]*flagDef),
	}
}

type flagSet struct {
	byName      map[string]*flagDef
	byShortName map[string]*flagDef

	args []string
}

// Args returns non-flag arguments.
func (fs *flagSet) Args() []string {
	return fs.args
}

// Has returns whether a flagset has a flag by a name.
func (fs *flagSet) Has(name string) bool {
	_, ok := fs.byName[name]
	return ok
}

// HasShort returns whether a flagset has a flag by a short name.
func (fs *flagSet) HasShort(name string) bool {
	_, ok := fs.byShortName[name]
	return ok
}

// Parse a set of command-line arguments as flags.
func (fs *flagSet) Parse(args []string) error {
	err := fs.parseFlags(args)
	if err != nil {
		return err
	}

	for _, f := range fs.byName {
		if err := fs.parseEnv(f); err != nil {
			return err
		}
	}

	return nil
}

func (fs *flagSet) parseFlags(args []string) error {
	for len(args) > 0 {
		arg := args[0]
		args = args[1:]

		switch {
		case arg == "--":
			fs.args = append(fs.args, args...)
			return nil
		case len(arg) < 2 || arg[0] != '-':
			fs.args = slices.Grow(fs.args, 1+len(args))
			fs.args = append(fs.args, arg)
			fs.args = append(fs.args, args...)
			return nil
		case arg[1] == '-':
			var err error
			args, err = fs.parseLong(arg, args)
			if err != nil {
				return err
			}
		default:
			var err error
			args, err = fs.parseShort(arg, args)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (fs *flagSet) parseLong(arg string, args []string) ([]string, error) {
	name, value, hasEqual := strings.Cut(arg[2:], "=")

	f, ok := fs.byName[name]
	if !ok {
		return nil, fmt.Errorf("unknown flag: --%s", name)
	}

	switch {
	case hasEqual:
	case f.boolVal != "":
		value = f.boolVal
	case len(args) > 0:
		value, args = args[0], args[1:]
	default:
		return nil, fmt.Errorf("missing argument for flag: --%s", name)
	}

	if err := f.set(value); err != nil {
		return nil, fmt.Errorf("invalid argument for flag --%s: %w", name, err)
	}

	return args, nil
}

func (fs *flagSet) parseShort(arg string, args []string) ([]string, error) {
	for i := 1; i < len(arg); i++ {
		short := string(arg[i])

		f, ok := fs.byShortName[short]
		if !ok {
			return nil, fmt.Errorf("unknown shorthand flag: '%s' in %s", short, arg)
		}

		isLastChar := i == len(arg)-1
		hasEqual := !isLastChar && arg[i+1] == '='
		hasMore := !isLastChar && arg[i+1] != '='

		var value string
		switch {
		case hasEqual:
			value = arg[i+2:]
			i = len(arg)
		case f.boolVal != "":
			value = f.boolVal
		case hasMore:
			value = arg[i+1:]
			i = len(arg)
		case len(args) > 0:
			value, args = args[0], args[1:]
		default:
			return nil, fmt.Errorf("missing argument for flag: '%s' in %s", short, arg)
		}

		if err := f.set(value); err != nil {
			return nil, fmt.Errorf("invalid argument for flag '%s' in %s: %w", short, arg, err)
		}
	}

	return args, nil
}

func (fs *flagSet) parseEnv(f *flagDef) error {
	if f.changed {
		return nil
	}

	for _, env := range f.env {
		v, ok := os.LookupEnv(env)
		if !ok {
			continue
		}

		if err := f.set(v); err != nil {
			return fmt.Errorf("invalid argument for env var %s: %w", env, err)
		}

		return nil
	}

	return nil
}
