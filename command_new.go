package clip

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/pflag"
)

// NewCommand creates a new command given a name and command options.
//
// By default, commands will print their help documentation when invoked.
// Different configuration options can be passed as a command is created, but
// the command returned will be immutable.
func NewCommand(name string, options ...func(*Command)) *Command {
	cmd := Command{
		name:          name,
		action:        printCommandHelp,
		subCommandMap: map[string]*Command{},
		writer:        os.Stdout,
		flagSet:       pflag.NewFlagSet(name, pflag.ContinueOnError),
	}

	cmd.flagSet.SetOutput(ioutil.Discard)

	// Overwrite defaults with passed options
	for _, o := range options {
		o(&cmd)
	}

	return &cmd
}

// AsHidden hides a command from documentation.
func AsHidden(cmd *Command) {
	cmd.hidden = true
}

// WithSummary adds a one-line description to a command.
func WithSummary(summary string) func(*Command) {
	return func(cmd *Command) {
		cmd.summary = summary
	}
}

// WithDescription adds a multi-line description to a command.
func WithDescription(description string) func(*Command) {
	return func(cmd *Command) {
		cmd.description = description
	}
}

// WithAction sets a Command's behavior when invoked.
func WithAction(action func(*Context) error) func(*Command) {
	return func(cmd *Command) {
		cmd.action = action
	}
}

// WithCommand adds a sub-command.
func WithCommand(subCmd *Command) func(*Command) {
	return func(cmd *Command) {
		if _, exists := cmd.subCommandMap[subCmd.Name()]; exists {
			panic(fmt.Sprintf("a sub-command with name %q already exists", subCmd.Name()))
		}
		cmd.subCommandMap[subCmd.Name()] = subCmd

		if !subCmd.hidden {
			cmd.visibleCommands = append(cmd.visibleCommands, subCmd)
		}
	}
}

// WithWriter sets the writer for writing output.
func WithWriter(writer io.Writer) func(*Command) {
	return func(cmd *Command) {
		cmd.writer = writer
	}
}
