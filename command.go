package clip

import "io"

// Command is a command or sub-command that can be run from the command-line.
//
// To create a new command with the default settings, use:
//  clip.NewCommand("command-name")
// rather than:
//  clip.Command{}
// The command type is immutable once created, so passing options to NewCommand
// is the only way to configure a command.
type Command struct {
	name        string
	description string
	action      func(*Context) error
	subCommands map[string]*Command
	writer      io.Writer
}

// Context is a command context with runtime metadata.
type Context struct {
	*Command

	Args []string
}

// Name is the name of the command.
func (cmd *Command) Name() string { return cmd.name }

// Description is a short description of the command.
func (cmd *Command) Description() string { return cmd.description }

// Run runs a command using a given set of args.
//
// The args parameter only includes the arguments passed specifically to a
// given command or sub-command.
func (cmd *Command) Run(args []string) error {
	if len(args) > 0 {
		if subCmd, ok := cmd.subCommands[args[0]]; ok {
			return subCmd.Run(args[1:])
		}
	}

	ctx := Context{
		Command: cmd,
		Args:    args,
	}

	return cmd.action(&ctx)
}
