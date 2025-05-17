package clip

import "io"

// Context is a command context with runtime metadata.
type Context struct {
	command *Command
	parent  *Context
}

// Name is the name of the command.
func (ctx *Context) Name() string {
	return ctx.command.Name()
}

// Summary is a one-line description of the command.
func (ctx *Context) Summary() string {
	return ctx.command.Summary()
}

// Description is a multi-line description of the command.
func (ctx *Context) Description() string {
	return ctx.command.Description()
}

// Writer is the writer for the command.
func (ctx *Context) Writer() io.Writer {
	for cur := ctx; cur != nil; cur = cur.parent {
		if cur.command.w != nil {
			return cur.command.w
		}
	}
	return defaultWriter
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

// Args returns the list of arguments.
func (ctx *Context) args() []string {
	return ctx.command.flagSet.Args()
}

// run runs the command with a given context.
func (ctx *Context) run(args []string) error {
	if err := ctx.command.flagSet.Parse(args[1:]); err != nil {
		return newUsageError(ctx, err.Error())
	}

	// Flag actions
	if wasSet, err := ctx.command.flagAction(ctx); wasSet {
		return err
	}

	// No sub commands or command action
	if len(ctx.command.subCommandMap) == 0 || len(ctx.args()) == 0 {
		return ctx.command.action(ctx)
	}

	// Sub commands, something passed
	subCmdName := ctx.args()[0]
	if subCmd, ok := ctx.command.subCommandMap[subCmdName]; ok {
		subCtx := Context{
			command: subCmd,
			parent:  ctx,
		}

		return subCtx.run(ctx.args())
	}

	return newUsageErrorf(ctx, "undefined sub-command: %s", subCmdName)
}
