package clip

import (
	"testing"

	"gotest.tools/assert"
)

func TestHelpContextFullName(t *testing.T) {
	var hctx *helpContext

	wasCalled := false
	action := func(ctx *Context) error {
		wasCalled = true
		hctx = newHelpContext(ctx)
		return nil
	}

	grandchild := NewCommand("grandchild", WithAction(action))
	child := NewCommand("child", WithCommand(grandchild))
	root := NewCommand("root", WithCommand(child))

	args := []string{root.Name(), child.Name(), grandchild.Name()}
	assert.NilError(t, root.Run(args))
	assert.Assert(t, wasCalled)
	assert.Check(t, hctx.FullName() == "root child grandchild")
}
