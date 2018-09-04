package clipflag

import "github.com/spf13/pflag"

type flag struct {
	name        string
	short       string
	summary     string
	description string
	define      func(*pflag.FlagSet)
	defineShort func(*pflag.FlagSet)

	// TODO: These
	envVar string
	hidden bool
}

func (f *flag) Name() string        { return f.name }
func (f *flag) Short() string       { return f.short }
func (f *flag) Summary() string     { return f.summary }
func (f *flag) Description() string { return f.description }
func (f *flag) Define(fs *pflag.FlagSet) {
	if f.short == "" {
		f.define(fs)
	} else {
		f.defineShort(fs)
	}
}
