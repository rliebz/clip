package clip

import (
	"errors"
	"fmt"
)

// Context is a command context with runtime metadata.
type Context struct {
	*Command

	parent *Context

	Args []string
}

// Parent is the context's parent context.
func (ctx *Context) Parent() *Context { return ctx.parent }

// Root is the context's root context.
func (ctx *Context) Root() *Context {
	cur := ctx
	for cur.parent != nil {
		cur = cur.parent
	}
	return cur
}

// run runs the command with a given context.
func (ctx *Context) run() error {
	// No sub commands
	if len(ctx.commands) == 0 {
		return ctx.action(ctx)
	}

	// Sub commands, but nothing passed
	if len(ctx.Args) == 0 {
		// TODO: Should help be printed?
		return errors.New("required sub-command not passed")
	}

	// Sub commands, something passed
	subCmdName := ctx.Args[0]
	if subCmd, ok := ctx.subCommandMap[subCmdName]; ok {
		subCtx := Context{
			Command: subCmd,
			Args:    ctx.Args[1:],
			parent:  ctx,
		}
		return subCtx.run()
	}

	return fmt.Errorf("undefined sub-command %q", subCmdName)

}
