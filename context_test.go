package clip

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
)

func TestContextParent(t *testing.T) {
	var pctx *Context

	wasCalled := false
	action := func(ctx *Context) error {
		wasCalled = true
		pctx = ctx.Parent()
		return nil
	}

	child := NewCommand(
		"foo",
		WithAction(action),
	)

	parent := NewCommand(
		"parent",
		WithCommand(child),
	)

	args := []string{parent.Name(), child.Name()}
	assert.NilError(t, parent.Run(args))
	assert.Assert(t, wasCalled)
	assert.Check(t, pctx.Name() == parent.Name())
	assert.Check(t, cmp.DeepEqual(pctx.Args(), args[1:]))
}

func TestContextParentNil(t *testing.T) {
	wasCalled := false
	action := func(ctx *Context) error {
		wasCalled = true
		assert.Check(t, ctx.Parent() == nil)
		return nil
	}

	cmd := NewCommand(
		"foo",
		WithAction(action),
	)

	assert.NilError(t, cmd.Run([]string{cmd.Name()}))
	assert.Assert(t, wasCalled)
}

func TestContextRoot(t *testing.T) {
	var tctx *Context

	action := func(ctx *Context) error {
		tctx = ctx
		return nil
	}

	foo := NewCommand("foo", WithAction(action))
	bar := NewCommand("bar", WithCommand(foo))
	baz := NewCommand("baz", WithCommand(bar))

	testCases := []struct {
		args []string
		cmd  *Command
	}{
		{[]string{"foo"}, foo},
		{[]string{"bar", "foo"}, bar},
		{[]string{"baz", "bar", "foo"}, baz},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("chain of %d command(s)", len(tt.args)), func(t *testing.T) {
			assert.NilError(t, tt.cmd.Run(tt.args))
			assert.Check(t, tctx.Root().Name() == tt.cmd.Name())
		})
	}
}
