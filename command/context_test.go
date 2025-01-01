package command

import (
	"fmt"
	"testing"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"
)

func TestContextParent(t *testing.T) {
	g := ghost.New(t)

	wasCalled := false
	var pctx *Context
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
	g.NoError(parent.Execute(args))
	g.Should(be.True(wasCalled))
	g.Should(be.Equal(pctx.Name(), parent.Name()))
	g.Should(be.DeepEqual(pctx.args(), args[1:]))
}

func TestContextParentNil(t *testing.T) {
	g := ghost.New(t)

	wasCalled := false
	action := func(ctx *Context) error {
		wasCalled = true
		g.Should(be.Nil(ctx.Parent()))
		return nil
	}

	cmd := New(
		"foo",
		WithAction(action),
	)

	g.NoError(cmd.Execute([]string{cmd.Name()}))
	g.Should(be.True(wasCalled))
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
			g := ghost.New(t)

			g.NoError(tt.cmd.Execute(tt.args))
			g.Should(be.Equal(tctx.Root().Name(), tt.cmd.Name()))
		})
	}
}
