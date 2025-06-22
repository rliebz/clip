package clip

import (
	"fmt"
	"io"
	"os"
)

// Command is a command or sub-command that can be run from the command-line.
//
// To create a new command with the default settings, use [NewCommand].
type Command struct {
	name        string
	summary     string
	description string
	hidden      bool
	action      func(*Context) error
	stdout      io.Writer
	stderr      io.Writer

	flagSet         *flagSet
	visibleCommands []*Command
	subCommandMap   map[string]*Command
	flagAction      func(*Context) (wasSet bool, err error)
}

// NewCommand creates a new command given a name and command options.
//
// By default, commands will print their help documentation when invoked.
func NewCommand(name string, options ...CommandOption) *Command {
	c := commandConfig{
		action: func(ctx *Context) error {
			return WriteHelp(ctx.Stdout(), ctx)
		},
		subCommandMap: map[string]*Command{},
		flagSet:       newFlagSet(),
		flagAction:    func(*Context) (bool, error) { return false, nil },
	}

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
		stdout:      c.stdout,
		stderr:      c.stderr,

		flagSet:         c.flagSet,
		visibleCommands: c.visibleCommands,
		subCommandMap:   c.subCommandMap,
		flagAction:      c.flagAction,
	}
}

// An CommandOption is used to commandConfigure a new command.
type CommandOption func(*commandConfig)

type commandConfig struct {
	summary     string
	description string
	hidden      bool
	action      func(*Context) error
	stdout      io.Writer
	stderr      io.Writer

	flagSet         *flagSet
	visibleCommands []*Command
	subCommandMap   map[string]*Command
	flagAction      func(*Context) (wasSet bool, err error)
}

// applyConditionalDefaults applies any conditionally-applied defaults.
// This includes things like a help or version flag that may not be applicable
// depending on which options are passed.
func applyConditionalDefaults(c *commandConfig) {
	if !c.flagSet.Has("help") {
		options := []FlagOption{
			FlagDescription("Print help and exit"),
			FlagAction(func(ctx *Context) error {
				return WriteHelp(ctx.Stdout(), ctx)
			}),
		}
		if !c.flagSet.HasShort("h") {
			options = append(options, FlagShort("h"))
		}

		ToggleFlag("help", options...)(c)
	}
}

// CommandHidden hides a command from documentation.
func CommandHidden(c *commandConfig) {
	c.hidden = true
}

// CommandSummary adds a one-line description to a command.
func CommandSummary(summary string) CommandOption {
	return func(c *commandConfig) {
		c.summary = summary
	}
}

// CommandDescription adds a multi-line description to a command.
func CommandDescription(description string) CommandOption {
	return func(c *commandConfig) {
		c.description = description
	}
}

// CommandAction sets a Command's behavior when invoked.
func CommandAction(action func(*Context) error) CommandOption {
	return func(c *commandConfig) {
		c.action = action
	}
}

// CommandSubCommand adds a sub-command.
func CommandSubCommand(subCmd *Command) CommandOption {
	return func(c *commandConfig) {
		if _, exists := c.subCommandMap[subCmd.Name()]; exists {
			panic(fmt.Sprintf("a sub-command with name %q already exists", subCmd.Name()))
		}
		c.subCommandMap[subCmd.Name()] = subCmd

		if !subCmd.hidden {
			c.visibleCommands = append(c.visibleCommands, subCmd)
		}
	}
}

// CommandStdout sets the writer for command output.
func CommandStdout(writer io.Writer) CommandOption {
	return func(c *commandConfig) {
		c.stdout = writer
	}
}

// CommandStderr sets the writer for command error output.
func CommandStderr(writer io.Writer) CommandOption {
	return func(c *commandConfig) {
		c.stderr = writer
	}
}

// addFlag registers a flag on a command.
//
// It is called after registering the flag on the command's flagset.
func (c *commandConfig) addFlag(f *flagDef) {
	c.flagSet.byName[f.name] = f
	if f.short != "" {
		c.flagSet.byShortName[f.short] = f
	}

	if f.action != nil {
		oldAction := c.flagAction
		c.flagAction = func(ctx *Context) (bool, error) {
			if wasSet, err := oldAction(ctx); wasSet {
				return true, err
			}
			if f.changed {
				return true, f.action(ctx)
			}
			return false, nil
		}
	}
}

// Name is the name of the command.
func (cmd *Command) Name() string { return cmd.name }

// Summary is a one-line description of the command.
func (cmd *Command) Summary() string { return cmd.summary }

// Description is a multi-line description of the command.
func (cmd *Command) Description() string { return cmd.description }

// Run runs a command.
//
// The args passed should begin with the name of the command itself.
// For the root command in most applications, the args will be [os.Args] and
// the result should be passed to [os.Exit].
func (cmd *Command) Run() int {
	ctx := &Context{
		command: cmd,
	}

	if err := ctx.run(os.Args); err != nil {
		ctx.printError(err)
		return exitCode(err)
	}

	return 0
}

// Execute runs a command using given args and returns the raw error.
//
// This function provides more fine-grained control than Run, and can be used
// in situations where handling arguments or errors needs more granular control.
func (cmd *Command) Execute(args []string) error {
	ctx := &Context{
		command: cmd,
	}

	if err := ctx.run(args); err != nil {
		return err
	}

	return nil
}
