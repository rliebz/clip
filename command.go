package clip

// Command is a command or sub-command that can be run from the command-line.
type Command struct {
	name   string
	action func(*Command) error
}

// Run runs a command using a given set of args.
func (cmd *Command) Run(args []string) error {
	if cmd.action != nil {
		return cmd.action(cmd)
	}

	return nil
}

type commandOption func(*Command)

// NewCommand creates a new command given a name and command options.
func NewCommand(name string, options ...commandOption) *Command {
	cmd := Command{name: name}
	for i := range options {
		options[i](&cmd)
	}

	return &cmd
}

// WithAction sets a Command's behavior when invoked.
func WithAction(action func(*Command) error) commandOption {
	return func(cmd *Command) {
		cmd.action = action
	}
}
