package command

import (
	"github.com/rliebz/clip"
)

// WithFlag adds a flag.
// Typically, flags from the clipflag package will be passed here.
func WithFlag(f clip.Flag) Option {
	return func(c *config) {
		f.Define(c.flagSet)
		if !f.Hidden() {
			c.visibleFlags = append(c.visibleFlags, f)
		}
	}
}

// WithActionFlag adds a flag that performs an action and nothing else.
// Flags such as --help or --version fall under this category.
//
// The action will occur if the flag is passed, regardless of the value, so
// typically clipflag.NewToggle will be used here.
func WithActionFlag(f clip.Flag, action func(*Context) error) Option {
	return func(c *config) {
		oldAction := c.flagAction
		f.Define(c.flagSet)
		if !f.Hidden() {
			c.visibleFlags = append(c.visibleFlags, f)
		}
		c.flagAction = func(ctx *Context) (bool, error) {
			if wasSet, err := oldAction(ctx); wasSet {
				return true, err
			}
			if c.flagSet.Changed(f.Name()) {
				return true, action(ctx)
			}
			return false, nil
		}
	}
}
