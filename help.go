package clip

import (
	"fmt"
	"io"
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

	maxFlagNameLen := 0
	for _, flag := range ctx.command.visibleFlags {
		if len(flag.Name()) > maxFlagNameLen {
			maxFlagNameLen = len(flag.Name())
		}
	}

	return &helpContext{
		Context:        ctx,
		maxCmdNameLen:  maxCmdNameLen,
		maxFlagNameLen: maxFlagNameLen,
	}
}

// helpContext is a wrapper around context with methods for printing help.
type helpContext struct {
	*Context

	maxCmdNameLen  int
	maxFlagNameLen int
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
func (ctx *helpContext) VisibleFlags() []*Flag { return ctx.command.visibleFlags }

const helpTemplateString = `{{.FullName}}{{if .Summary}} - {{.Summary}}{{end}}{{if .Description}}

{{.Description}}{{end}}{{if .VisibleCommands}}

Commands:{{range .VisibleCommands}}
  {{padCommand .Name}}{{if .Summary}}{{.Summary}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

Flags:{{range .VisibleFlags}}
  {{printFlagShort .Short}}--{{padFlag .Name}}{{if .Summary}}{{.Summary}}{{end}}{{end}}{{end}}
`

var printCommandHelp = func(ctx *Context) error {
	return writeCommandHelp(ctx.Writer(), ctx)
}

func writeCommandHelp(wr io.Writer, ctx *Context) error {
	hctx := newHelpContext(ctx)
	t := template.New("help").Funcs(template.FuncMap{
		"padCommand":     getCommandPadder(hctx),
		"padFlag":        getFlagPadder(hctx),
		"printFlagShort": printFlagShort,
	})
	t = template.Must(t.Parse(helpTemplateString))
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

func getFlagPadder(ctx *helpContext) func(string) string {
	s := fmt.Sprintf("%%-%ds", ctx.maxFlagNameLen+2)
	return func(text string) string {
		return fmt.Sprintf(s, text)
	}
}
