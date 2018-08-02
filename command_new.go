package clip

import (
	"fmt"
	"io"
	"os"
)

type commandOption func(*Command)

// NewCommand creates a new command given a name and command options.
func NewCommand(name string, options ...commandOption) *Command {
	cmd := Command{
		name:        name,
		action:      printCommandHelp,
		subCommands: map[string]*Command{},
		writer:      os.Stdout,
	}

	// Overwrite defaults with passed options
	for i := range options {
		options[i](&cmd)
	}

	return &cmd
}

// WithAction sets a Command's behavior when invoked.
func WithAction(action func(*Context) error) commandOption {
	return func(cmd *Command) {
		cmd.action = action
	}
}

// WithDescription adds a short description to a command.
func WithDescription(description string) commandOption {
	return func(cmd *Command) {
		cmd.description = description
	}
}

// WithWriter sets the writer for writing output.
func WithWriter(writer io.Writer) commandOption {
	return func(cmd *Command) {
		cmd.writer = writer
	}
}

// WithCommand adds a sub-command.
func WithCommand(subCmd *Command) commandOption {
	return func(cmd *Command) {
		if _, exists := cmd.subCommands[subCmd.Name()]; exists {
			panic(fmt.Sprintf("a sub-command with name %q already exists", subCmd.Name()))
		}
		cmd.subCommands[subCmd.Name()] = subCmd
	}
}
