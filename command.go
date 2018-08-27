package clip

import (
	"errors"
	"io"
)

// Command is a command or sub-command that can be run from the command-line.
//
// To create a new command with the default settings, use:
//  clip.NewCommand("command-name")
// rather than:
//  clip.Command{}
// The command type is immutable once created, so passing options to NewCommand
// is the only way to configure a command.
type Command struct {
	name          string
	summary       string
	description   string
	action        func(*Context) error
	commands      []*Command
	subCommandMap map[string]*Command
	writer        io.Writer
}

// Name is the name of the command.
func (cmd *Command) Name() string { return cmd.name }

// Summary is a one-line description of the command.
func (cmd *Command) Summary() string { return cmd.summary }

// Description is a multi-line description of the command.
func (cmd *Command) Description() string { return cmd.description }

// Commands is the list of sub-commands in order.
func (cmd *Command) Commands() []*Command { return cmd.commands }

// Run runs a command using a given set of args.
//
// The args passed should begin with the name of the command itself.
// For the root command in most applications, the args will be os.Args.
func (cmd *Command) Run(args []string) error {
	if len(args) == 0 {
		return errors.New("no arguments were passed")
	}

	ctx := &Context{
		Command: cmd,
		args:    args[1:],
	}

	if err := ctx.run(); err != nil {
		printErrorHelp(err)
		return err
	}

	return nil
}
