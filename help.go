package clip

import (
	"html/template"
)

const helpTemplateString = `{{.Name}}{{if .Summary}} - {{.Summary}}{{end}}{{if .Description}}

{{.Description}}{{end}}{{if .Commands}}

Commands:{{range .Commands}}
  {{.Name}}{{if .Summary}} - {{.Summary}}{{end}}{{end}}{{end}}
`

func printCommandHelp(ctx *Context) error {
	t := template.Must(template.New("help").Parse(helpTemplateString))
	return t.Execute(ctx.writer, ctx)
}
