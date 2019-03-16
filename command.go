package clip

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/spf13/pflag"
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
	name        string
	summary     string
	description string
	hidden      bool
	action      func(*Context) error
	writer      io.Writer

	flagSet         *pflag.FlagSet
	visibleCommands []*Command
	subCommandMap   map[string]*Command
	flagAction      func(*Context) (wasSet bool, err error)
}

// Name is the name of the command.
func (cmd *Command) Name() string { return cmd.name }

// Summary is a one-line description of the command.
func (cmd *Command) Summary() string { return cmd.summary }

// Description is a multi-line description of the command.
func (cmd *Command) Description() string { return cmd.description }

// Writer is the writer for the command.
func (cmd *Command) Writer() io.Writer { return cmd.writer }

// VisibleCommands is the list of sub-commands in order.
func (cmd *Command) VisibleCommands() []*Command { return cmd.visibleCommands }

// Execute runs a command using given args and returns the raw error.
//
// This function provides more fine-grained control than Run, and can be used
// in situations where handling arguments or errors needs more granular control.
func (cmd *Command) Execute(args []string) error {
	ctx := &Context{
		Command: cmd,
	}

	if len(args) == 0 {
		return errors.New("no arguments were passed")
	}

	if err := ctx.run(args); err != nil {
		return err
	}

	return nil
}

// Run runs a command.
//
// The args passed should begin with the name of the command itself.
// For the root command in most applications, the args will be os.Args.
func (cmd *Command) Run() int {
	if err := cmd.Execute(os.Args); err != nil {
		l := log.New(cmd.Writer(), "", 0)
		printError(l, err)
		return getExitCode(err)
	}

	return 0
}
