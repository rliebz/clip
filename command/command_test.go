package command

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/assert/cmp"

	"github.com/rliebz/clip"
	"github.com/rliebz/clip/flag"
)

func TestCommandName(t *testing.T) {
	name := "foo"
	command := New("foo")
	assert.Equal(t, command.Name(), name)
}

func TestCommandSummary(t *testing.T) {
	summary := "some summary"
	command := New("foo", WithSummary(summary))
	assert.Equal(t, command.Summary(), summary)
}

func TestCommandDescription(t *testing.T) {
	description := "some description"
	command := New("foo", WithDescription(description))
	assert.Equal(t, command.Description(), description)
}

func TestCommandWriter(t *testing.T) {
	writer := new(bytes.Buffer)
	command := New("foo", WithWriter(writer))
	assert.Equal(t, command.writer, writer)
}

func TestCommandExecuteHelp(t *testing.T) {
	cmdName := "foo"
	output := new(bytes.Buffer)
	command := New(
		cmdName,
		WithSummary("some summary"),
		WithDescription("some description"),
		WithWriter(output),
	)
	assert.NilError(t, command.Execute([]string{cmdName}))
	helpText := output.String()
	assert.Check(t, cmp.Contains(helpText, command.Name()))
	assert.Check(t, cmp.Contains(helpText, command.Summary()))
	assert.Check(t, cmp.Contains(helpText, command.Description()))
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
		t.Run(
			fmt.Sprintf("Flag %q/%q passed %s", tt.flag.Name(), tt.flag.Short(), tt.passed),
			func(t *testing.T) {
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
					assert.NilError(t, err)
					assert.Check(t, !flagActionCalled)
					assert.Check(t, cmp.Contains(helpText, command.Name()))
				case callsAction:
					assert.NilError(t, err)
					assert.Check(t, flagActionCalled)
					assert.Check(t, helpText == "")
				case hasError:
					assert.Error(t, err, "unknown shorthand flag: 'h' in -h")
				}
			})
	}
}

func TestCommandAction(t *testing.T) {
	wasCalled := false
	action := func(*Context) error {
		wasCalled = true
		return nil
	}

	command := New("foo", WithAction(action))

	assert.Check(t, !wasCalled)
	assert.NilError(t, command.Execute([]string{command.Name()}))
	assert.Check(t, wasCalled)
}

func TestCommandActionError(t *testing.T) {
	err := errors.New("some error")

	wasCalled := false
	action := func(*Context) error {
		wasCalled = true
		return err
	}

	command := New("foo", WithAction(action))

	assert.Check(t, !wasCalled)
	assert.Error(t, command.Execute([]string{command.Name()}), err.Error())
	assert.Check(t, wasCalled)
}

func TestCommandFlagAction(t *testing.T) {
	wasCalled := false
	action := func(*Context) error {
		wasCalled = true
		return nil
	}
	flagValue := false
	flg := flag.NewBool(&flagValue, "my-flag")

	command := New("foo", WithActionFlag(flg, action))

	assert.Check(t, !wasCalled)
	assert.Check(t, !flagValue)
	assert.NilError(t, command.Execute([]string{command.Name(), "--my-flag"}))
	assert.Check(t, wasCalled)
	assert.Check(t, flagValue)
}

func TestCommandFlagCorrectAction(t *testing.T) {
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

	assert.Check(t, !wasCalled)
	assert.NilError(t, command.Execute(
		[]string{command.Name(), "--my-flag", "--second-flag", "bar"}),
	)
	assert.Check(t, wasCalled)
	assert.Check(t, !wrongWasCalled)
}

func TestCommandFlagActionError(t *testing.T) {
	err := errors.New("some error")

	wasCalled := false
	action := func(*Context) error {
		wasCalled = true
		return err
	}

	flagValue := false
	f := flag.NewBool(&flagValue, "my-flag")

	command := New("foo", WithActionFlag(f, action))

	assert.Check(t, !wasCalled)
	assert.Check(t, !flagValue)
	assert.Error(t, command.Execute([]string{command.Name(), "--my-flag"}), err.Error())
	assert.Check(t, wasCalled)
	assert.Check(t, flagValue)
}

func TestCommandArgs(t *testing.T) {
	cmdName := "foo"
	args := []string{"a", "b", "c"}

	wasCalled := false
	action := func(ctx *Context) error {
		wasCalled = true
		assert.Check(t, cmp.DeepEqual(args, ctx.args()))
		return nil
	}

	command := New(
		cmdName,
		WithAction(action),
	)

	cliArgs := append([]string{cmdName}, args...)
	assert.NilError(t, command.Execute(cliArgs))
	assert.Assert(t, wasCalled)
}

func TestSubCommandArgs(t *testing.T) {
	cmdName := "foo"
	subCmdName := "bar"
	args := []string{"a", "b", "c"}

	subCmdWasCalled := false
	subCmdAction := func(ctx *Context) error {
		subCmdWasCalled = true
		assert.Check(t, cmp.DeepEqual(args, ctx.args()))
		return nil
	}
	defer func() { assert.Check(t, subCmdWasCalled) }()

	subCommand := New(
		subCmdName,
		WithAction(subCmdAction),
	)

	cmdWasCalled := false
	cmdAction := func(*Context) error {
		cmdWasCalled = true
		return nil
	}
	defer func() { assert.Check(t, !cmdWasCalled) }()

	command := New(
		cmdName,
		WithAction(cmdAction),
		WithCommand(subCommand),
	)

	cliArgs := append([]string{cmdName, subCmdName}, args...)
	assert.NilError(t, command.Execute(cliArgs))
}

func TestSubCommandDuplicates(t *testing.T) {
	assert.Assert(t, cmp.Panics(func() {
		New(
			"foo",
			WithCommand(New("bar")),
			WithCommand(New("bar")),
		)
	}))
}

func TestCommandNoArgs(t *testing.T) {
	command := New("foo")
	assert.Error(t, command.Execute([]string{}), "no arguments were passed")
}

func TestCommandNoSubCommands(t *testing.T) {
	cmdName := "my-command"

	defer func(original func(ctx *Context) error) {
		printCommandHelp = original
	}(printCommandHelp)

	helpPrinted := false
	printCommandHelp = func(ctx *Context) error {
		helpPrinted = true
		assert.Check(t, ctx.Name() == cmdName)
		return nil
	}

	child := New("foo")
	parent := New(cmdName, WithCommand(child))

	args := []string{parent.Name()}
	assert.NilError(t, parent.Execute(args))
	assert.Assert(t, helpPrinted)
}

func TestCommandNonExistentSubCommand(t *testing.T) {
	command := New("foo")
	parent := New("bar", WithCommand(command))
	assert.Error(
		t,
		parent.Execute([]string{parent.Name(), "wrong"}),
		"undefined sub-command: wrong",
	)
}

func TestRun(t *testing.T) {
	defer func(args []string) { os.Args = args }(os.Args)
	os.Args = []string{"foo"}
	buf := new(bytes.Buffer)
	command := New(
		"foo",
		WithAction(func(*Context) error { return nil }),
		WithWriter(buf),
	)

	assert.Assert(t, command.Run() == 0)
	assert.Check(t, cmp.DeepEqual(buf.String(), ""))
}

func TestRunError(t *testing.T) {
	defer func(args []string) { os.Args = args }(os.Args)
	os.Args = []string{"foo"}
	err := errors.New("oops")
	buf := new(bytes.Buffer)
	command := New(
		"foo",
		WithAction(func(*Context) error { return err }),
		WithWriter(buf),
	)

	assert.Check(t, command.Run() == 1)
	assert.Check(t, cmp.Contains(buf.String(), err.Error()))
}

func TestRunExitError(t *testing.T) {
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
	assert.Assert(t, errors.As(err, &exitErr))
	assert.Check(t, command.Run() == exitErr.ExitCode())
	assert.Check(t, cmp.Contains(buf.String(), err.Error()))
}

func TestRunUsageError(t *testing.T) {
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

	assert.Check(t, command.Run() == 1)
	assert.Check(t, cmp.Contains(buf.String(), errMessage))
}
