package command

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

	child := New(
		"foo",
		WithAction(action),
	)

	parent := New(
		"parent",
		WithCommand(child),
	)

	args := []string{parent.Name(), child.Name()}
	assert.NilError(t, parent.Execute(args))
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

	cmd := New(
		"foo",
		WithAction(action),
	)

	assert.NilError(t, cmd.Execute([]string{cmd.Name()}))
	assert.Assert(t, wasCalled)
}

func TestContextRoot(t *testing.T) {
	var tctx *Context

	action := func(ctx *Context) error {
		tctx = ctx
		return nil
	}

	foo := New("foo", WithAction(action))
	bar := New("bar", WithCommand(foo))
	baz := New("baz", WithCommand(bar))

	tests := []struct {
		args []string
		cmd  *Command
	}{
		{[]string{"foo"}, foo},
		{[]string{"bar", "foo"}, bar},
		{[]string{"baz", "bar", "foo"}, baz},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("chain of %d command(s)", len(tt.args)), func(t *testing.T) {
			assert.NilError(t, tt.cmd.Execute(tt.args))
			assert.Check(t, tctx.Root().Name() == tt.cmd.Name())
		})
	}
}
