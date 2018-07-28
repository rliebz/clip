package clip

import (
	"bytes"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
)

func TestCommandName(t *testing.T) {
	description := "some description"
	command := NewCommand("foo", WithDescription(description))
	assert.Equal(t, command.Description(), description)
}

func TestCommandRunHelp(t *testing.T) {
	output := new(bytes.Buffer)
	command := NewCommand(
		"foo",
		WithDescription("some description"),
		WithWriter(output),
	)
	assert.NilError(t, command.Run([]string{}))
	helpText := output.String()
	assert.Check(t, cmp.Contains(helpText, command.Name()))
	assert.Check(t, cmp.Contains(helpText, command.Description()))
}

func TestCommandAction(t *testing.T) {
	wasCalled := false
	action := func(ctx *Context) error {
		wasCalled = true
		return nil
	}

	command := NewCommand(
		"foo",
		WithAction(action),
	)

	assert.Check(t, !wasCalled)
	assert.NilError(t, command.Run([]string{}))
	assert.Check(t, wasCalled)
}

func TestCommandArgs(t *testing.T) {
	args := []string{"a", "b", "c"}

	wasCalled := false
	action := func(ctx *Context) error {
		wasCalled = true
		assert.DeepEqual(t, ctx.Args, args)
		return nil
	}

	command := NewCommand(
		"foo",
		WithAction(action),
	)

	assert.NilError(t, command.Run(args))
	assert.Assert(t, wasCalled)
}
