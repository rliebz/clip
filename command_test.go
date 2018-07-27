package clip

import (
	"testing"

	"gotest.tools/assert"
)

func TestCommandRunDefault(t *testing.T) {
	command := NewCommand("foo")
	command.Run([]string{})
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
	command.Run([]string{})
	assert.Check(t, wasCalled)
}
