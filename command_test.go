package clip

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
)

func TestCommandName(t *testing.T) {
	name := "foo"
	command := NewCommand("foo")
	assert.Equal(t, command.Name(), name)
}

func TestCommandSummary(t *testing.T) {
	summary := "some summary"
	command := NewCommand("foo", WithSummary(summary))
	assert.Equal(t, command.Summary(), summary)
}

func TestCommandDescription(t *testing.T) {
	description := "some description"
	command := NewCommand("foo", WithDescription(description))
	assert.Equal(t, command.Description(), description)
}

func TestCommandWriter(t *testing.T) {
	writer := new(bytes.Buffer)
	command := NewCommand("foo", WithWriter(writer))
	assert.Equal(t, command.Writer(), writer)
}

func TestCommandExecuteHelp(t *testing.T) {
	cmdName := "foo"
	output := new(bytes.Buffer)
	command := NewCommand(
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

func TestCommandAction(t *testing.T) {
	wasCalled := false
	action := func(ctx *Context) error {
		wasCalled = true
		return nil
	}

	command := NewCommand("foo", WithAction(action))

	assert.Check(t, !wasCalled)
	assert.NilError(t, command.Execute([]string{command.Name()}))
	assert.Check(t, wasCalled)
}

func TestCommandActionError(t *testing.T) {
	err := errors.New("some error")

	wasCalled := false
	action := func(ctx *Context) error {
		wasCalled = true
		return err
	}

	command := NewCommand("foo", WithAction(action))

	assert.Check(t, !wasCalled)
	assert.Error(t, command.Execute([]string{command.Name()}), err.Error())
	assert.Check(t, wasCalled)
}

func TestCommandArgs(t *testing.T) {
	cmdName := "foo"
	args := []string{"a", "b", "c"}

	wasCalled := false
	action := func(ctx *Context) error {
		wasCalled = true
		assert.DeepEqual(t, ctx.args, args)
		return nil
	}

	command := NewCommand(
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
		assert.DeepEqual(t, ctx.args, args)
		return nil
	}
	defer func() { assert.Check(t, subCmdWasCalled) }()

	subCommand := NewCommand(
		subCmdName,
		WithAction(subCmdAction),
	)

	cmdWasCalled := false
	cmdAction := func(ctx *Context) error {
		cmdWasCalled = true
		return nil
	}
	defer func() { assert.Check(t, !cmdWasCalled) }()

	command := NewCommand(
		cmdName,
		WithAction(cmdAction),
		WithCommand(subCommand),
	)

	cliArgs := append([]string{cmdName, subCmdName}, args...)
	assert.NilError(t, command.Execute(cliArgs))
}

func TestSubCommandDuplicates(t *testing.T) {
	assert.Assert(t, cmp.Panics(func() {
		NewCommand(
			"foo",
			WithCommand(NewCommand("bar")),
			WithCommand(NewCommand("bar")),
		)
	}))
}

func TestCommandNoArgs(t *testing.T) {
	command := NewCommand("foo")
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

	child := NewCommand("foo")
	parent := NewCommand(cmdName, WithCommand(child))

	args := []string{parent.Name()}
	assert.NilError(t, parent.Execute(args))
	assert.Assert(t, helpPrinted)
}

func TestCommandNonExistentSubCommand(t *testing.T) {
	command := NewCommand("foo")
	parent := NewCommand("bar", WithCommand(command))
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
	command := NewCommand(
		"foo",
		WithAction(func(ctx *Context) error { return nil }),
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
	command := NewCommand(
		"foo",
		WithAction(func(ctx *Context) error { return err }),
		WithWriter(buf),
	)

	assert.Check(t, command.Run() == 1)
	assert.Check(t, cmp.Contains(buf.String(), err.Error()))
}

func TestRunExitError(t *testing.T) {
	defer func(args []string) { os.Args = args }(os.Args)
	os.Args = []string{"foo"}
	err := NewError("oops", 3).(exitError)
	buf := new(bytes.Buffer)
	command := NewCommand(
		"foo",
		WithAction(func(ctx *Context) error { return err }),
		WithWriter(buf),
	)

	assert.Check(t, command.Run() == err.ExitCode())
	assert.Check(t, cmp.Contains(buf.String(), err.Error()))
}

func TestRunUsageError(t *testing.T) {
	defer func(args []string) { os.Args = args }(os.Args)
	os.Args = []string{"foo"}
	errMessage := "oops"
	buf := new(bytes.Buffer)
	command := NewCommand(
		"foo",
		WithAction(func(ctx *Context) error {
			return newUsageError(ctx, errMessage)
		}),
		WithWriter(buf),
	)

	assert.Check(t, command.Run() == 1)
	assert.Check(t, cmp.Contains(buf.String(), errMessage))
}
