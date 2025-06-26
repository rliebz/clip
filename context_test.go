package clip

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

	parent := NewCommand(
		"parent",
		SubCommand(
			"child",
			CommandAction(action),
		),
	)

	args := []string{"parent", "child"}
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

	cmd := NewCommand(
		"foo",
		CommandAction(action),
	)

	g.NoError(cmd.Execute([]string{cmd.Name()}))
	g.Should(be.True(wasCalled))
}

func TestContextRoot(t *testing.T) {
	tests := []struct {
		args      []string
		fooCalled bool
		barCalled bool
		bazCalled bool
	}{
		{
			args:      []string{"foo"},
			fooCalled: true,
		},
		{
			args:      []string{"foo", "bar"},
			barCalled: true,
		},
		{
			args:      []string{"foo", "bar", "baz"},
			bazCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("chain of %d command(s)", len(tt.args)), func(t *testing.T) {
			g := ghost.New(t)

			fooCalled := false
			barCalled := false
			bazCalled := false

			cmd := NewCommand(
				"foo",
				CommandAction(func(ctx *Context) error {
					fooCalled = true
					g.Should(be.Equal(ctx.Root().Name(), "foo"))
					return nil
				}),
				SubCommand(
					"bar",
					CommandAction(func(ctx *Context) error {
						barCalled = true
						g.Should(be.Equal(ctx.Root().Name(), "foo"))
						return nil
					}),
					SubCommand(
						"baz",
						CommandAction(func(ctx *Context) error {
							bazCalled = true
							g.Should(be.Equal(ctx.Root().Name(), "foo"))
							return nil
						}),
					),
				),
			)

			g.NoError(cmd.Execute(tt.args))

			g.Should(be.Equal(fooCalled, tt.fooCalled))
			g.Should(be.Equal(barCalled, tt.barCalled))
			g.Should(be.Equal(bazCalled, tt.bazCalled))
		})
	}
}
