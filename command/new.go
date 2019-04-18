package command

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/rliebz/clip/clipflag"
	"github.com/spf13/pflag"
)

// New creates a new command given a name and command options.
//
// By default, commands will print their help documentation when invoked.
// Different configuration options can be passed as a command is created, but
// the command returned will be immutable.
func New(name string, options ...Option) *Command {
	c := config{
		action:        printCommandHelp,
		subCommandMap: map[string]*Command{},
		writer:        os.Stdout,
		flagSet:       pflag.NewFlagSet(name, pflag.ContinueOnError),
		flagAction:    func(ctx *Context) (bool, error) { return false, nil },
	}

	c.flagSet.SetOutput(ioutil.Discard)

	// Overwrite defaults with passed options
	for _, o := range options {
		o(&c)
	}

	applyConditionalDefaults(&c)

	return &Command{
		name:        name,
		summary:     c.summary,
		description: c.description,
		hidden:      c.hidden,
		action:      c.action,
		writer:      c.writer,

		flagSet:         c.flagSet,
		visibleCommands: c.visibleCommands,
		visibleFlags:    c.visibleFlags,
		subCommandMap:   c.subCommandMap,
		flagAction:      c.flagAction,
	}
}

// An Option is used to configure a new command.
type Option func(*config)

type config struct {
	summary     string
	description string
	hidden      bool
	action      func(*Context) error
	writer      io.Writer

	flagSet         *pflag.FlagSet
	visibleCommands []*Command
	visibleFlags    []Flag
	subCommandMap   map[string]*Command
	flagAction      func(*Context) (wasSet bool, err error)
}

// applyConditionalDefaults applies any conditionally-applied defaults.
// This includes things like a help or version flag that may not be applicable
// depending on which options are passed.
func applyConditionalDefaults(c *config) {
	if c.flagSet.Lookup("help") == nil {
		options := []clipflag.Option{
			clipflag.WithSummary("Print help and exit"),
		}
		if c.flagSet.ShorthandLookup("h") == nil {
			options = append(options, clipflag.WithShort("h"))
		}

		flag := clipflag.NewToggle("help", options...)
		WithActionFlag(flag, printCommandHelp)(c)
	}
}

// AsHidden hides a command from documentation.
func AsHidden(c *config) {
	c.hidden = true
}

// WithSummary adds a one-line description to a command.
func WithSummary(summary string) Option {
	return func(c *config) {
		c.summary = summary
	}
}

// WithDescription adds a multi-line description to a command.
func WithDescription(description string) Option {
	return func(c *config) {
		c.description = description
	}
}

// WithAction sets a Command's behavior when invoked.
func WithAction(action func(*Context) error) Option {
	return func(c *config) {
		c.action = action
	}
}

// WithCommand adds a sub-command.
func WithCommand(subCmd *Command) Option {
	return func(c *config) {
		if _, exists := c.subCommandMap[subCmd.Name()]; exists {
			panic(fmt.Sprintf("a sub-command with name %q already exists", subCmd.Name()))
		}
		c.subCommandMap[subCmd.Name()] = subCmd

		if !subCmd.hidden {
			c.visibleCommands = append(c.visibleCommands, subCmd)
		}
	}
}

// WithWriter sets the writer for writing output.
func WithWriter(writer io.Writer) Option {
	return func(c *config) {
		c.writer = writer
	}
}
