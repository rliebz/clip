package clip

import "io"

// Command is a command or sub-command that can be run from the command-line.
type Command struct {
	// Command definition
	name        string
	description string
	action      func(*Command) error
	writer      io.Writer

	// Runtime metadata
	args []string
}

// Name is the name of the command.
func (cmd *Command) Name() string { return cmd.name }

// Description is a description of the command.
func (cmd *Command) Description() string { return cmd.description }

// Run runs a command using a given set of args.
func (cmd *Command) Run(args []string) error {
	cmd.args = args

	return cmd.action(cmd)
}

// Args returns a list of arguments passed to the command when run.
func (cmd *Command) Args() []string {
	return cmd.args
}
