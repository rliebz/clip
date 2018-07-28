package clip

import (
	"html/template"
)

const helpTemplateString = `{{.Name}} - {{.Description}}`

func printCommandHelp(cmd *Command) error {
	t := template.Must(template.New("help").Parse(helpTemplateString))
	return t.Execute(cmd.writer, cmd)
}
