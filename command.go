package clip

import "io"

// Command is a command or sub-command that can be run from the command-line.
type Command struct {
	// Command definition
	name        string
	description string
	action      func(*Context) error
	writer      io.Writer
}

// Context is a command context with runtime metadata.
type Context struct {
	*Command

	Args []string
}

// Name is the name of the command.
func (cmd *Command) Name() string { return cmd.name }

// Description is a description of the command.
func (cmd *Command) Description() string { return cmd.description }

// Run runs a command using a given set of args.
func (cmd *Command) Run(args []string) error {
	ctx := Context{
		Command: cmd,
		Args:    args,
	}

	return cmd.action(&ctx)
}
