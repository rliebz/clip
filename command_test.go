package clip

import (
	"bytes"
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

func TestCommandRunHelp(t *testing.T) {
	cmdName := "foo"
	output := new(bytes.Buffer)
	command := NewCommand(
		cmdName,
		WithSummary("some summary"),
		WithDescription("some description"),
		WithWriter(output),
	)
	assert.NilError(t, command.Run([]string{cmdName}))
	helpText := output.String()
	assert.Check(t, cmp.Contains(helpText, command.Name()))
	assert.Check(t, cmp.Contains(helpText, command.Summary()))
	assert.Check(t, cmp.Contains(helpText, command.Description()))
}

func TestCommandAction(t *testing.T) {
	cmdName := "foo"

	wasCalled := false
	action := func(ctx *Context) error {
		wasCalled = true
		return nil
	}

	command := NewCommand(
		cmdName,
		WithAction(action),
	)

	assert.Check(t, !wasCalled)
	assert.NilError(t, command.Run([]string{cmdName}))
	assert.Check(t, wasCalled)
}

func TestCommandArgs(t *testing.T) {
	cmdName := "foo"
	args := []string{"a", "b", "c"}

	wasCalled := false
	action := func(ctx *Context) error {
		wasCalled = true
		assert.DeepEqual(t, ctx.Args, args)
		return nil
	}

	command := NewCommand(
		cmdName,
		WithAction(action),
	)

	cliArgs := append([]string{cmdName}, args...)
	assert.NilError(t, command.Run(cliArgs))
	assert.Assert(t, wasCalled)
}

func TestSubCommandArgs(t *testing.T) {
	cmdName := "foo"
	subCmdName := "bar"
	args := []string{"a", "b", "c"}

	subCmdWasCalled := false
	subCmdAction := func(ctx *Context) error {
		subCmdWasCalled = true
		assert.DeepEqual(t, ctx.Args, args)
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
	assert.NilError(t, command.Run(cliArgs))
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
	assert.Error(t, command.Run([]string{}), "no arguments were passed")
}

func TestCommandNoSubCommands(t *testing.T) {
	command := NewCommand("foo")
	parent := NewCommand("bar", WithCommand(command))
	assert.Error(t, parent.Run([]string{parent.Name()}), "required sub-command not passed")
}

func TestCommandNonExistentSubCommand(t *testing.T) {
	command := NewCommand("foo")
	parent := NewCommand("bar", WithCommand(command))
	assert.Error(
		t,
		parent.Run([]string{parent.Name(), "wrong"}),
		`undefined sub-command "wrong"`,
	)
}
