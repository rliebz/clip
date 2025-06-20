package clip

import (
	"bytes"
	"testing"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"
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

	grandchild := NewCommand("grandchild", CommandAction(action))
	child := NewCommand("child", CommandSubCommand(grandchild))
	root := NewCommand("root", CommandSubCommand(child))

	args := []string{root.Name(), child.Name(), grandchild.Name()}
	g.NoError(root.Execute(args))
	g.Should(be.True(wasCalled))
	g.Should(be.Equal(hctx.FullName(), "root child grandchild"))
}

func TestHelpCommands(t *testing.T) {
	g := ghost.New(t)

	buf := new(bytes.Buffer)
	root := NewCommand(
		"root",
		CommandStdout(buf),
		CommandSubCommand(NewCommand("child-one", CommandSummary("1"))),
		CommandSubCommand(NewCommand("child-two", CommandSummary("2"))),
		CommandSubCommand(NewCommand("child-three", CommandSummary("3"))),
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
	root := NewCommand(
		"root",
		CommandStdout(buf),
		CommandSubCommand(NewCommand("visible")),
		CommandSubCommand(NewCommand("hidden", CommandHidden)),
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
	root := NewCommand(
		"root",
		CommandStdout(buf),
		ToggleFlag("visible"),
		ToggleFlag("hidden", FlagHidden),
	)

	args := []string{root.Name()}
	g.NoError(root.Execute(args))

	output := buf.String()
	g.Should(be.StringContaining(output, "visible"))
	g.ShouldNot(be.StringContaining(output, "hidden"))
}

func Test_printCommandHelp_placeholder(t *testing.T) {
	g := ghost.New(t)

	buf := new(bytes.Buffer)
	root := NewCommand(
		"root",
		CommandStdout(buf),
		StringFlag(new(string), "default-string"),
		StringFlag(new(string), "override-string", FlagPlaceholder("somevalue")),
		BoolFlag(new(bool), "default-bool"),
		BoolFlag(new(bool), "override-bool", FlagPlaceholder("somebool")),
	)

	args := []string{root.Name()}
	g.NoError(root.Execute(args))

	output := buf.String()
	g.Should(be.StringContaining(output, "--default-string <string>"))
	g.Should(be.StringContaining(output, "--override-string <somevalue>"))
	g.Should(be.StringContaining(output, "--default-bool\n"))
	g.Should(be.StringContaining(output, "--override-bool=<somebool>"))
}
