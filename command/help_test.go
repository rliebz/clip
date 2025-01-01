package command

import (
	"bytes"
	"testing"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"

	"github.com/rliebz/clip/flag"
)

func TestHelpContextFullName(t *testing.T) {
	g := ghost.New(t)

	wasCalled := false
	var hctx *helpContext
	action := func(ctx *Context) error {
		wasCalled = true
		hctx = newHelpContext(ctx)
		return nil
	}

	grandchild := New("grandchild", WithAction(action))
	child := New("child", WithCommand(grandchild))
	root := New("root", WithCommand(child))

	args := []string{root.Name(), child.Name(), grandchild.Name()}
	g.NoError(root.Execute(args))
	g.Should(be.True(wasCalled))
	g.Should(be.Equal(hctx.FullName(), "root child grandchild"))
}

func TestHelpCommands(t *testing.T) {
	g := ghost.New(t)

	buf := new(bytes.Buffer)
	root := New(
		"root",
		WithWriter(buf),
		WithCommand(New("child-one", WithSummary("1"))),
		WithCommand(New("child-two", WithSummary("2"))),
		WithCommand(New("child-three", WithSummary("3"))),
	)

	args := []string{root.Name()}
	g.NoError(root.Execute(args))

	output := buf.String()
	g.Should(be.StringContaining(output, "child-one    1"))
	g.Should(be.StringContaining(output, "child-two    2"))
	g.Should(be.StringContaining(output, "child-three  3"))
}

func TestHidden(t *testing.T) {
	g := ghost.New(t)

	buf := new(bytes.Buffer)
	root := New(
		"root",
		WithWriter(buf),
		WithCommand(New("visible")),
		WithCommand(New("hidden", AsHidden)),
	)

	args := []string{root.Name()}
	g.NoError(root.Execute(args))

	output := buf.String()
	g.Should(be.StringContaining(output, "visible"))
	g.ShouldNot(be.StringContaining(output, "hidden"))
}

func TestHiddenFlags(t *testing.T) {
	g := ghost.New(t)

	buf := new(bytes.Buffer)
	root := New(
		"root",
		WithWriter(buf),
		WithFlag(flag.NewToggle("visible")),
		WithFlag(flag.NewToggle("hidden", flag.AsHidden)),
	)

	args := []string{root.Name()}
	g.NoError(root.Execute(args))

	output := buf.String()
	g.Should(be.StringContaining(output, "visible"))
	g.ShouldNot(be.StringContaining(output, "hidden"))
}
