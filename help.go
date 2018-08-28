package clip

import (
	"fmt"
	"html/template"
	"io"
)

// newHelpContext creates a helpContext from a Context.
func newHelpContext(ctx *Context) *helpContext {
	max := 0
	for _, cmd := range ctx.commands {
		if len(cmd.Name()) > max {
			max = len(cmd.Name())
		}
	}

	return &helpContext{
		Context:       ctx,
		maxCmdNameLen: max,
	}
}

// helpContext is a wrapper around context with methods for printing help.
type helpContext struct {
	*Context

	maxCmdNameLen int
}

func (ctx *helpContext) FullName() string {
	name := ctx.Name()

	cur := ctx.Parent()
	for cur != nil {
		name = fmt.Sprintf("%s %s", cur.Name(), name)
		cur = cur.Parent()
	}

	return name
}

const helpTemplateString = `{{.FullName}}{{if .Summary}} - {{.Summary}}{{end}}{{if .Description}}

{{.Description}}{{end}}{{if .Commands}}

Commands:{{range .Commands}}
  {{padCommand .Name}}{{if .Summary}}{{.Summary}}{{end}}{{end}}{{end}}
`

var printCommandHelp = func(ctx *Context) error {
	return writeCommandHelp(ctx.writer, ctx)
}

func writeCommandHelp(wr io.Writer, ctx *Context) error {
	hctx := newHelpContext(ctx)
	t := template.New("help").Funcs(template.FuncMap{
		"padCommand": getCommandPadder(hctx),
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
