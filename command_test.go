package clip

import (
	"testing"

	"gotest.tools/assert"
)

func TestCommandRunDefault(t *testing.T) {
	command := NewCommand("foo")
	assert.NilError(t, command.Run([]string{}))
}

func TestCommandAction(t *testing.T) {
	wasCalled := false
	action := func(cmd *Command) error {
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
	action := func(cmd *Command) error {
		wasCalled = true
		assert.DeepEqual(t, cmd.Args(), args)
		return nil
	}

	command := NewCommand(
		"foo",
		WithAction(action),
	)

	assert.NilError(t, command.Run(args))
	assert.Assert(t, wasCalled)
}
