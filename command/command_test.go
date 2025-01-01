package command

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"

	"github.com/rliebz/clip"
	"github.com/rliebz/clip/flag"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name            string
		opts            []Option
		wantDescription string
		wantSummary     string
	}{
		{
			name: "defaults",
		},
		{
			name: "WithSummary",
			opts: []Option{
				WithSummary("some summary"),
			},
			wantSummary: "some summary",
		},
		{
			name: "WithDescription",
			opts: []Option{
				WithDescription("some description"),
			},
			wantDescription: "some description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := ghost.New(t)

			command := New("foo", tt.opts...)
			g.Should(be.Equal(command.Name(), "foo"))
			g.Should(be.Equal(command.Description(), tt.wantDescription))
			g.Should(be.Equal(command.Summary(), tt.wantSummary))
		})
	}
}

func TestCommandWriter(t *testing.T) {
	g := ghost.New(t)

	writer := new(bytes.Buffer)
	command := New("foo", WithWriter(writer))
	g.Should(be.Equal[io.Writer](command.writer, writer))
}

func TestCommandExecuteHelp(t *testing.T) {
	g := ghost.New(t)

	cmdName := "foo"
	output := new(bytes.Buffer)
	command := New(
		cmdName,
		WithSummary("some summary"),
		WithDescription("some description"),
		WithWriter(output),
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
		flag     clip.Flag
		passed   string
		behavior int
	}{
		{flag.NewToggle("help"), "--help", callsAction},
		{flag.NewToggle("help"), "-h", hasError},
		{flag.NewToggle("help", flag.WithShort("h")), "--help", callsAction},
		{flag.NewToggle("help", flag.WithShort("h")), "-h", callsAction},
		{flag.NewToggle("help", flag.WithShort("n")), "--help", callsAction},
		{flag.NewToggle("help", flag.WithShort("n")), "-h", hasError},
		{flag.NewToggle("not-help"), "--help", callsHelp},
		{flag.NewToggle("not-help"), "-h", callsHelp},
		{flag.NewToggle("not-help", flag.WithShort("h")), "--help", callsHelp},
		{flag.NewToggle("not-help", flag.WithShort("h")), "-h", callsAction},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("Flag %q/%q passed %s", tt.flag.Name(), tt.flag.Short(), tt.passed)
		t.Run(name, func(t *testing.T) {
			g := ghost.New(t)

			cmdName := "foo"
			output := new(bytes.Buffer)
			flagActionCalled := false
			command := New(
				cmdName,
				WithActionFlag(
					tt.flag,
					func(*Context) error {
						flagActionCalled = true
						return nil
					},
				),
				WithWriter(output),
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

	command := New("foo", WithAction(action))
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

	command := New("foo", WithAction(action))
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
	flg := flag.NewBool(&flagValue, "my-flag")

	command := New("foo", WithActionFlag(flg, action))
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
	notCalledFlag := flag.NewBool(&notCalledValue, "not-called")
	correctValue := false
	correctFlag := flag.NewBool(&correctValue, "my-flag")
	secondValue := false
	secondFlag := flag.NewBool(&secondValue, "second-flag")

	subCommand := New("bar", WithAction(wrongAction))

	command := New(
		"foo",
		WithCommand(subCommand),
		WithActionFlag(notCalledFlag, wrongAction),
		WithActionFlag(correctFlag, action),
		WithActionFlag(secondFlag, wrongAction),
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
	f := flag.NewBool(&flagValue, "my-flag")

	command := New("foo", WithActionFlag(f, action))

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

	command := New(
		cmdName,
		WithAction(action),
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

	subCommand := New(
		subCmdName,
		WithAction(subCmdAction),
	)

	cmdWasCalled := false
	cmdAction := func(*Context) error {
		cmdWasCalled = true
		return nil
	}
	defer func() { g.Should(be.False(cmdWasCalled)) }()

	command := New(
		cmdName,
		WithAction(cmdAction),
		WithCommand(subCommand),
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

	New(
		"foo",
		WithCommand(New("bar")),
		WithCommand(New("bar")),
	)
}

func TestCommandNoArgs(t *testing.T) {
	g := ghost.New(t)

	command := New("foo")
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

	child := New("foo")
	parent := New(cmdName, WithCommand(child))

	args := []string{parent.Name()}
	err := parent.Execute(args)
	g.NoError(err)

	g.Should(be.True(helpPrinted))
}

func TestCommandNonExistentSubCommand(t *testing.T) {
	g := ghost.New(t)

	command := New("foo")
	parent := New("bar", WithCommand(command))
	err := parent.Execute([]string{parent.Name(), "wrong"})
	g.Should(be.ErrorEqual(err, "undefined sub-command: wrong"))
}

func TestRun(t *testing.T) {
	g := ghost.New(t)

	defer func(args []string) { os.Args = args }(os.Args)
	os.Args = []string{"foo"}

	buf := new(bytes.Buffer)
	command := New(
		"foo",
		WithAction(func(*Context) error { return nil }),
		WithWriter(buf),
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
	command := New(
		"foo",
		WithAction(func(*Context) error { return wantErr }),
		WithWriter(buf),
	)

	g.Should(be.Equal(command.Run(), 1))
	g.Should(be.StringContaining(buf.String(), wantErr.Error()))
}

func TestRunExitError(t *testing.T) {
	g := ghost.New(t)

	defer func(args []string) { os.Args = args }(os.Args)
	os.Args = []string{"foo"}

	err := NewError("oops", 3)
	buf := new(bytes.Buffer)
	command := New(
		"foo",
		WithAction(func(*Context) error { return err }),
		WithWriter(buf),
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
	command := New(
		"foo",
		WithAction(func(ctx *Context) error {
			return newUsageError(ctx, errMessage)
		}),
		WithWriter(buf),
	)

	g.Should(be.Equal(command.Run(), 1))
	g.Should(be.StringContaining(buf.String(), errMessage))
}
