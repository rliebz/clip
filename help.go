package clip

import (
	"html/template"
)

const helpTemplateString = `{{.Name}} - {{.Description}}`

func printCommandHelp(ctx *Context) error {
	t := template.Must(template.New("help").Parse(helpTemplateString))
	return t.Execute(ctx.writer, ctx)
}
