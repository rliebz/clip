package clip

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name            string
		opts            []CommandOption
		wantDescription string
		wantSummary     string
	}{
		{
			name: "defaults",
		},
		{
			name: "WithSummary",
			opts: []CommandOption{
				CommandSummary("some summary"),
			},
			wantSummary: "some summary",
		},
		{
			name: "WithDescription",
			opts: []CommandOption{
				CommandDescription("some description"),
			},
			wantDescription: "some description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := ghost.New(t)

			command := NewCommand("foo", tt.opts...)
			g.Should(be.Equal(command.Name(), "foo"))
			g.Should(be.Equal(command.Description(), tt.wantDescription))
			g.Should(be.Equal(command.Summary(), tt.wantSummary))
		})
	}
}

func TestCommandWriter(t *testing.T) {
	g := ghost.New(t)

	writer := new(bytes.Buffer)
	command := NewCommand("foo", CommandWriter(writer))
	g.Should(be.Equal[io.Writer](command.writer(), writer))
}

func TestCommandExecuteHelp(t *testing.T) {
	g := ghost.New(t)

	cmdName := "foo"
	output := new(bytes.Buffer)
	command := NewCommand(
		cmdName,
		CommandSummary("some summary"),
		CommandDescription("some description"),
		CommandWriter(output),
	)

	err := command.Execute([]string{cmdName})
	g.NoError(err)

	helpText := output.String()
	g.Should(be.StringContaining(helpText, command.Name()))
	g.Should(be.StringContaining(helpText, command.Summary()))
	g.Should(be.StringContaining(helpText, command.Description()))
}

func TestCommandDefaultHelpFlag(t *testing.T) {
	const (
		callsHelp   = iota
		callsAction = iota
		hasError    = iota
	)

	tests := []struct {
		flag     Flag
		passed   string
		behavior int
	}{
		{NewToggleFlag("help"), "--help", callsAction},
		{NewToggleFlag("help"), "-h", hasError},
		{NewToggleFlag("help", FlagShort("h")), "--help", callsAction},
		{NewToggleFlag("help", FlagShort("h")), "-h", callsAction},
		{NewToggleFlag("help", FlagShort("n")), "--help", callsAction},
		{NewToggleFlag("help", FlagShort("n")), "-h", hasError},
		{NewToggleFlag("not-help"), "--help", callsHelp},
		{NewToggleFlag("not-help"), "-h", callsHelp},
		{NewToggleFlag("not-help", FlagShort("h")), "--help", callsHelp},
		{NewToggleFlag("not-help", FlagShort("h")), "-h", callsAction},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("Flag %q/%q passed %s", tt.flag.Name(), tt.flag.Short(), tt.passed)
		t.Run(name, func(t *testing.T) {
			g := ghost.New(t)

			cmdName := "foo"
			output := new(bytes.Buffer)
			flagActionCalled := false
			command := NewCommand(
				cmdName,
				CommandActionFlag(
					tt.flag,
					func(*Context) error {
						flagActionCalled = true
						return nil
					},
				),
				CommandWriter(output),
			)
			err := command.Execute([]string{cmdName, tt.passed})
			helpText := output.String()
			switch tt.behavior {
			case callsHelp:
				g.NoError(err)
				g.Should(be.False(flagActionCalled))
				g.Should(be.StringContaining(helpText, command.Name()))
			case callsAction:
				g.NoError(err)
				g.Should(be.True(flagActionCalled))
				g.Should(be.Zero(helpText))
			case hasError:
				g.Should(be.ErrorEqual(err, "unknown shorthand flag: 'h' in -h"))
			}
		})
	}
}

func TestCommandAction(t *testing.T) {
	g := ghost.New(t)

	wasCalled := false
	action := func(*Context) error {
		wasCalled = true
		return nil
	}

	command := NewCommand("foo", CommandAction(action))
	g.Should(be.False(wasCalled))

	err := command.Execute([]string{command.Name()})
	g.NoError(err)

	g.Should(be.True(wasCalled))
}

func TestCommandActionError(t *testing.T) {
	g := ghost.New(t)

	wantErr := errors.New("some error")

	wasCalled := false
	action := func(*Context) error {
		wasCalled = true
		return wantErr
	}

	command := NewCommand("foo", CommandAction(action))
	g.Should(be.False(wasCalled))

	err := command.Execute([]string{command.Name()})
	g.Should(be.Equal(err, wantErr))
	g.Should(be.True(wasCalled))
}

func TestCommandFlagAction(t *testing.T) {
	g := ghost.New(t)

	wasCalled := false
	action := func(*Context) error {
		wasCalled = true
		return nil
	}
	flagValue := false
	flg := NewBoolFlag(&flagValue, "my-flag")

	command := NewCommand("foo", CommandActionFlag(flg, action))
	g.Should(be.False(wasCalled))
	g.Should(be.False(flagValue))

	err := command.Execute([]string{command.Name(), "--my-flag"})
	g.NoError(err)
	g.Should(be.True(wasCalled))
	g.Should(be.True(flagValue))
}

func TestCommandFlagCorrectAction(t *testing.T) {
	g := ghost.New(t)

	wasCalled := false
	wrongWasCalled := false
	action := func(*Context) error {
		wasCalled = true
		return nil
	}
	wrongAction := func(*Context) error {
		wrongWasCalled = true
		return nil
	}

	notCalledValue := false
	notCalledFlag := NewBoolFlag(&notCalledValue, "not-called")
	correctValue := false
	correctFlag := NewBoolFlag(&correctValue, "my-flag")
	secondValue := false
	secondFlag := NewBoolFlag(&secondValue, "second-flag")

	subCommand := NewCommand("bar", CommandAction(wrongAction))

	command := NewCommand(
		"foo",
		CommandSubCommand(subCommand),
		CommandActionFlag(notCalledFlag, wrongAction),
		CommandActionFlag(correctFlag, action),
		CommandActionFlag(secondFlag, wrongAction),
	)

	g.Should(be.False(wasCalled))

	err := command.Execute([]string{command.Name(), "--my-flag", "--second-flag", "bar"})
	g.NoError(err)

	g.Should(be.True(wasCalled))
	g.Should(be.False(wrongWasCalled))
}

func TestCommandFlagActionError(t *testing.T) {
	g := ghost.New(t)

	wantErr := errors.New("some error")

	wasCalled := false
	action := func(*Context) error {
		wasCalled = true
		return wantErr
	}

	flagValue := false
	f := NewBoolFlag(&flagValue, "my-flag")

	command := NewCommand("foo", CommandActionFlag(f, action))

	g.Should(be.False(wasCalled))
	g.Should(be.False(flagValue))

	err := command.Execute([]string{command.Name(), "--my-flag"})
	g.Should(be.Equal(err, wantErr))

	g.Should(be.True(wasCalled))
	g.Should(be.True(flagValue))
}

func TestCommandArgs(t *testing.T) {
	g := ghost.New(t)

	cmdName := "foo"
	args := []string{"a", "b", "c"}

	wasCalled := false
	action := func(ctx *Context) error {
		wasCalled = true
		g.Should(be.DeepEqual(ctx.args(), args))
		return nil
	}

	command := NewCommand(
		cmdName,
		CommandAction(action),
	)

	cliArgs := append([]string{cmdName}, args...)
	err := command.Execute(cliArgs)
	g.NoError(err)
	g.Should(be.True(wasCalled))
}

func TestSubCommandArgs(t *testing.T) {
	g := ghost.New(t)

	cmdName := "foo"
	subCmdName := "bar"
	args := []string{"a", "b", "c"}

	subCmdWasCalled := false
	subCmdAction := func(ctx *Context) error {
		subCmdWasCalled = true
		g.Should(be.DeepEqual(ctx.args(), args))
		return nil
	}
	defer func() { g.Should(be.True(subCmdWasCalled)) }()

	subCommand := NewCommand(
		subCmdName,
		CommandAction(subCmdAction),
	)

	cmdWasCalled := false
	cmdAction := func(*Context) error {
		cmdWasCalled = true
		return nil
	}
	defer func() { g.Should(be.False(cmdWasCalled)) }()

	command := NewCommand(
		cmdName,
		CommandAction(cmdAction),
		CommandSubCommand(subCommand),
	)

	cliArgs := append([]string{cmdName, subCmdName}, args...)
	err := command.Execute(cliArgs)
	g.NoError(err)
}

func TestSubCommandDuplicates(t *testing.T) {
	g := ghost.New(t)

	defer func() {
		g.Should(be.Equal(recover(), `a sub-command with name "bar" already exists`))
	}()

	NewCommand(
		"foo",
		CommandSubCommand(NewCommand("bar")),
		CommandSubCommand(NewCommand("bar")),
	)
}

func TestCommandNoArgs(t *testing.T) {
	g := ghost.New(t)

	command := NewCommand("foo")
	err := command.Execute([]string{})
	g.Should(be.ErrorEqual(err, "no arguments were passed"))
}

func TestCommandNoSubCommands(t *testing.T) {
	g := ghost.New(t)

	cmdName := "my-command"

	defer func(original func(ctx *Context) error) {
		printCommandHelp = original
	}(printCommandHelp)

	helpPrinted := false
	printCommandHelp = func(ctx *Context) error {
		helpPrinted = true
		g.Should(be.Equal(ctx.Name(), cmdName))
		return nil
	}

	child := NewCommand("foo")
	parent := NewCommand(cmdName, CommandSubCommand(child))

	args := []string{parent.Name()}
	err := parent.Execute(args)
	g.NoError(err)

	g.Should(be.True(helpPrinted))
}

func TestCommandNonExistentSubCommand(t *testing.T) {
	g := ghost.New(t)

	command := NewCommand("foo")
	parent := NewCommand("bar", CommandSubCommand(command))
	err := parent.Execute([]string{parent.Name(), "wrong"})
	g.Should(be.ErrorEqual(err, "undefined sub-command: wrong"))
}

func TestRun(t *testing.T) {
	g := ghost.New(t)

	defer func(args []string) { os.Args = args }(os.Args)
	os.Args = []string{"foo"}

	buf := new(bytes.Buffer)
	command := NewCommand(
		"foo",
		CommandAction(func(*Context) error { return nil }),
		CommandWriter(buf),
	)

	g.Should(be.Zero(command.Run()))
	g.Should(be.Zero(buf.String()))
}

func TestRunError(t *testing.T) {
	g := ghost.New(t)

	defer func(args []string) { os.Args = args }(os.Args)
	os.Args = []string{"foo"}

	wantErr := errors.New("oops")
	buf := new(bytes.Buffer)
	command := NewCommand(
		"foo",
		CommandAction(func(*Context) error { return wantErr }),
		CommandWriter(buf),
	)

	g.Should(be.Equal(command.Run(), 1))
	g.Should(be.StringContaining(buf.String(), wantErr.Error()))
}

func TestRunExitError(t *testing.T) {
	g := ghost.New(t)

	defer func(args []string) { os.Args = args }(os.Args)
	os.Args = []string{"foo"}

	err := NewExitError("oops", 3)
	buf := new(bytes.Buffer)
	command := NewCommand(
		"foo",
		CommandAction(func(*Context) error { return err }),
		CommandWriter(buf),
	)

	var exitErr exitError
	g.Must(be.ErrorAs(err, &exitErr))
	g.Should(be.Equal(command.Run(), exitErr.ExitCode()))
	g.Should(be.StringContaining(buf.String(), err.Error()))
}

func TestRunUsageError(t *testing.T) {
	g := ghost.New(t)

	defer func(args []string) { os.Args = args }(os.Args)
	os.Args = []string{"foo"}

	errMessage := "oops"
	buf := new(bytes.Buffer)
	command := NewCommand(
		"foo",
		CommandAction(func(ctx *Context) error {
			return newUsageError(ctx, errMessage)
		}),
		CommandWriter(buf),
	)

	g.Should(be.Equal(command.Run(), 1))
	g.Should(be.StringContaining(buf.String(), errMessage))
}

func TestParse(t *testing.T) {
	tests := []struct {
		args         []string
		childSFlag   string
		parentSFlag  string
		childFlag    bool
		parentFlag   bool
		childCalled  bool
		parentCalled bool
	}{
		{
			args:         []string{"foo"},
			parentCalled: true,
		},
		{
			args:         []string{"foo", "-f"},
			parentFlag:   true,
			parentCalled: true,
		},
		{
			args:         []string{"foo", "--flag"},
			parentFlag:   true,
			parentCalled: true,
		},
		{
			args:         []string{"foo", "--sflag", "bar"},
			parentSFlag:  "bar",
			parentCalled: true,
		},
		{
			args:         []string{"foo", "-fs", "bar"},
			parentFlag:   true,
			parentSFlag:  "bar",
			parentCalled: true,
		},
		{
			args:         []string{"foo", "--sflag", "bar", "--flag"},
			parentFlag:   true,
			parentSFlag:  "bar",
			parentCalled: true,
		},
		{
			args:         []string{"foo", "--flag=true"},
			parentFlag:   true,
			parentCalled: true,
		},
		{
			args:        []string{"foo", "child"},
			childCalled: true,
		},
		{
			args:        []string{"foo", "child", "--flag"},
			childFlag:   true,
			childCalled: true,
		},
		{
			args:        []string{"foo", "--flag", "child"},
			childCalled: true,
			parentFlag:  true,
		},
		{
			args:        []string{"foo", "--flag=true", "child"},
			childCalled: true,
			parentFlag:  true,
		},
		{
			args:        []string{"foo", "--flag", "child", "--flag"},
			childFlag:   true,
			childCalled: true,
			parentFlag:  true,
		},
		{
			args:        []string{"foo", "-fs=bar", "child"},
			childCalled: true,
			parentFlag:  true,
			parentSFlag: "bar",
		},
		{
			args:        []string{"foo", "-f", "-s=bar", "child"},
			childCalled: true,
			parentFlag:  true,
			parentSFlag: "bar",
		},
		{
			args:        []string{"foo", "--sflag", "bar", "child"},
			childCalled: true,
			parentSFlag: "bar",
		},
		{
			args:        []string{"foo", "-fs", "bar", "child", "-fs", "baz"},
			childFlag:   true,
			childSFlag:  "baz",
			childCalled: true,
			parentFlag:  true,
			parentSFlag: "bar",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("args: %v", tt.args), func(t *testing.T) {
			g := ghost.New(t)

			childFlag := false
			childSFlag := ""
			childCalled := false
			child := NewCommand(
				"child",
				CommandFlag(NewStringFlag(&childSFlag, "sflag", FlagShort("s"))),
				CommandFlag(NewBoolFlag(&childFlag, "flag", FlagShort("f"))),
				CommandAction(func(*Context) error {
					childCalled = true
					return nil
				}),
			)

			parentFlag := false
			parentSFlag := ""
			parentCalled := false
			cmd := NewCommand(
				"foo",
				CommandFlag(NewStringFlag(&parentSFlag, "sflag", FlagShort("s"))),
				CommandFlag(NewBoolFlag(&parentFlag, "flag", FlagShort("f"))),
				CommandSubCommand(child),
				CommandAction(func(*Context) error {
					parentCalled = true
					return nil
				}),
			)

			g.NoError(cmd.Execute(tt.args))

			g.Should(be.Equal(childCalled, tt.childCalled))
			g.Should(be.Equal(childFlag, tt.childFlag))
			g.Should(be.Equal(childSFlag, tt.childSFlag))
			g.Should(be.Equal(parentCalled, tt.parentCalled))
			g.Should(be.Equal(parentFlag, tt.parentFlag))
			g.Should(be.Equal(parentSFlag, tt.parentSFlag))
		})
	}
}

func TestParseError(t *testing.T) {
	tests := []struct {
		args []string
		err  string
	}{
		{
			args: []string{"foo", "--bad"},
			err:  "unknown flag: bad",
		},
		{
			args: []string{"foo", "-b"},
			err:  "unknown shorthand flag: 'b' in -b",
		},
		{
			args: []string{"foo", "-bad"},
			err:  "unknown shorthand flag: 'd' in -bad",
		},
		{
			args: []string{"foo", "bad"},
			err:  "undefined sub-command: bad",
		},
		{
			args: []string{"foo", "child", "--bad"},
			err:  "unknown flag: bad",
		},
		{
			args: []string{"foo", "child", "-b"},
			err:  "unknown shorthand flag: 'b' in -b",
		},
		{
			args: []string{"foo", "child", "-bad"},
			err:  "unknown shorthand flag: 'd' in -bad",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("args: %v", tt.args), func(t *testing.T) {
			g := ghost.New(t)

			childFlag := false
			childCalled := false
			child := NewCommand(
				"child",
				CommandFlag(NewBoolFlag(&childFlag, "flag")),
				CommandAction(func(*Context) error {
					childCalled = true
					return nil
				}),
			)

			parentFlag := false
			parentCalled := false
			cmd := NewCommand(
				"foo",
				CommandFlag(NewBoolFlag(&parentFlag, "flag")),
				CommandSubCommand(child),
				CommandAction(func(*Context) error {
					parentCalled = true
					return nil
				}),
			)

			g.Should(be.ErrorEqual(cmd.Execute(tt.args), tt.err))
			g.Should(be.False(childCalled))
			g.Should(be.False(parentCalled))
		})
	}
}
