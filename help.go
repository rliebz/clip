package clip

import (
	"html/template"
)

const helpTemplateString = `{{.Name}}{{if .Summary}} - {{.Summary}}{{end}}{{if .Description}}

{{.Description}}{{end}}
`

func printCommandHelp(ctx *Context) error {
	t := template.Must(template.New("help").Parse(helpTemplateString))
	return t.Execute(ctx.writer, ctx)
}
