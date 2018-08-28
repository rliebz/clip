package clip

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
)

// newUsageError creates an error which causes help to be printed.
func newUsageError(ctx *Context, message string) usageError { // nolint: unparam
	return usageError{
		context: ctx,
		message: message,
	}
}

// newUsageErrorf creates an error which causes help to be printed.
func newUsageErrorf(ctx *Context, format string, a ...interface{}) usageError { // nolint: unparam
	return usageError{
		context: ctx,
		message: fmt.Sprintf(format, a...),
	}
}

// usageError is an error caused by incorrect usage.
type usageError struct {
	context *Context
	message string
}

func (e usageError) Error() string { return e.message }
func (e usageError) ErrorContext() string {
	b := new(bytes.Buffer)
	_ = writeCommandHelp(b, e.context)
	return b.String()
}

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

func printCommandHelp(ctx *Context) error {
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
