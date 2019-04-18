package command

import (
	"bytes"
	"strings"
	"testing"

	"gotest.tools/assert"
	"gotest.tools/assert/cmp"

	"github.com/rliebz/clip/clipflag"
)

func TestHelpContextFullName(t *testing.T) {
	var hctx *helpContext

	wasCalled := false
	action := func(ctx *Context) error {
		wasCalled = true
		hctx = newHelpContext(ctx)
		return nil
	}

	grandchild := New("grandchild", WithAction(action))
	child := New("child", WithCommand(grandchild))
	root := New("root", WithCommand(child))

	args := []string{root.Name(), child.Name(), grandchild.Name()}
	assert.NilError(t, root.Execute(args))
	assert.Assert(t, wasCalled)
	assert.Check(t, hctx.FullName() == "root child grandchild")
}

func TestHelpCommands(t *testing.T) {
	buf := new(bytes.Buffer)
	root := New(
		"root",
		WithWriter(buf),
		WithCommand(New("child-one", WithSummary("1"))),
		WithCommand(New("child-two", WithSummary("2"))),
		WithCommand(New("child-three", WithSummary("3"))),
	)

	args := []string{root.Name()}
	assert.NilError(t, root.Execute(args))

	output := buf.String()
	assert.Check(t, cmp.Contains(output, "child-one    1"))
	assert.Check(t, cmp.Contains(output, "child-two    2"))
	assert.Check(t, cmp.Contains(output, "child-three  3"))
}

func TestHidden(t *testing.T) {
	buf := new(bytes.Buffer)
	root := New(
		"root",
		WithWriter(buf),
		WithCommand(New("visible")),
		WithCommand(New("hidden", AsHidden)),
	)

	args := []string{root.Name()}
	assert.NilError(t, root.Execute(args))

	output := buf.String()
	assert.Check(t, cmp.Contains(output, "visible"))
	assert.Check(t, !strings.Contains(output, "hidden"))
}

func TestHiddenFlags(t *testing.T) {
	buf := new(bytes.Buffer)
	root := New(
		"root",
		WithWriter(buf),
		WithFlag(clipflag.NewToggle("visible")),
		WithFlag(clipflag.NewToggle("hidden", clipflag.AsHidden)),
	)

	args := []string{root.Name()}
	assert.NilError(t, root.Execute(args))

	output := buf.String()
	assert.Check(t, cmp.Contains(output, "visible"))
	assert.Check(t, !strings.Contains(output, "hidden"))
}
