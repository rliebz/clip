package clip

// Context is a command context with runtime metadata.
type Context struct {
	*Command

	args   []string
	parent *Context
}

// Args returns the list of arguments.
func (ctx *Context) Args() []string { return ctx.args }

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
func (ctx *Context) run(args []string) error {
	if err := parse(ctx, args[1:]); err != nil {
		return newUsageError(ctx, err.Error())
	}
	ctx.args = ctx.flagSet.Args()

	// No sub commands or command action
	if len(ctx.subCommandMap) == 0 || len(ctx.args) == 0 {
		return ctx.action(ctx)
	}

	// Sub commands, something passed
	subCmdName := ctx.args[0]
	if subCmd, ok := ctx.subCommandMap[subCmdName]; ok {
		subCtx := Context{
			Command: subCmd,
			parent:  ctx,
		}
		return subCtx.run(ctx.args)
	}

	return newUsageErrorf(ctx, "undefined sub-command: %s", subCmdName)
}
