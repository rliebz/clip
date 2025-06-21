package clip

import (
	_ "embed"
	"fmt"
	"io"
	"iter"
	"strings"
	"text/template"
)

// newHelpContext creates a helpContext from a Context.
func newHelpContext(ctx *Context) *helpContext {
	maxCmdNameLen := 0
	for _, cmd := range ctx.command.visibleCommands {
		if len(cmd.Name()) > maxCmdNameLen {
			maxCmdNameLen = len(cmd.Name())
		}
	}

	return &helpContext{
		Context:       ctx,
		maxCmdNameLen: maxCmdNameLen,
	}
}

// helpContext is a wrapper around context with methods for printing help.
type helpContext struct {
	*Context

	maxCmdNameLen int
}

func (ctx *helpContext) FullName() string {
	name := ctx.command.Name()

	cur := ctx.Parent()
	for cur != nil {
		name = fmt.Sprintf("%s %s", cur.command.Name(), name)
		cur = cur.Parent()
	}

	return name
}

// VisibleCommands is the list of sub-commands in order.
func (ctx *helpContext) VisibleCommands() []*Command { return ctx.command.visibleCommands }

// VisibleFlags is the list of flags in order.
func (ctx *helpContext) VisibleFlags() iter.Seq[*flagDef] {
	return func(yield func(*flagDef) bool) {
		for _, flag := range ctx.command.flagSet.byName {
			if flag.Hidden() {
				continue
			}
			if !yield(flag) {
				return
			}
		}
	}
}

// TODO: Default values
//
//go:embed help.tmpl
var helpTemplate string

var printCommandHelp = func(ctx *Context) error {
	return writeCommandHelp(ctx.Stdout(), ctx)
}

func writeCommandHelp(wr io.Writer, ctx *Context) error {
	hctx := newHelpContext(ctx)
	t := template.New("help").Funcs(template.FuncMap{
		"join":           stringsJoin,
		"pad":            pad,
		"padCommand":     getCommandPadder(hctx),
		"printFlagShort": printFlagShort,
	})
	t = template.Must(t.Parse(helpTemplate))
	return t.Execute(wr, hctx)
}

func getCommandPadder(ctx *helpContext) func(string) string {
	s := fmt.Sprintf("%%-%ds", ctx.maxCmdNameLen+2)
	return func(text string) string {
		return fmt.Sprintf(s, text)
	}
}

func printFlagShort(short string) string {
	if short == "" {
		return "    "
	}
	return fmt.Sprintf("-%s, ", short)
}

// stringsJoin reverse the args for [strings.Join] to work better with templates.
func stringsJoin(sep string, s []string) string {
	return strings.Join(s, sep)
}

// pad a string with a number of leading spaces on each new line.
func pad(size int, text string) string {
	padding := strings.Repeat(" ", size)
	return padding + strings.ReplaceAll(
		strings.TrimSpace(text),
		"\n",
		"\n"+padding,
	)
}
