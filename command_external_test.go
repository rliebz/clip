package clip_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/rliebz/ghost"
	"github.com/rliebz/ghost/be"

	"github.com/rliebz/clip"
)

func TestNewCommand_help(t *testing.T) {
	g := ghost.New(t)

	var loud bool
	name := "world"
	boolTrue := true

	buf := new(bytes.Buffer)
	hello := clip.NewCommand(
		"hello",
		clip.CommandStdout(buf),
		clip.CommandSummary("Greet the world"),
		clip.CommandDescription(`This is a command that will say hello to
a person, or to the world.`),
		clip.BoolFlag(
			&loud,
			"loud",
			clip.FlagDescription(`Whether to pump up the volume to max.
This is a very long multi-line, very multi-line
description of things to come and more!
`),
			clip.FlagDeprecated("Please don't be loud."),
			clip.FlagEnv("HELLO_LOUD", "LOUD"),
			clip.FlagShort("l"),
		),
		clip.BoolFlag(
			&boolTrue,
			"bool-a",
		),
		clip.BoolFlag(
			&boolTrue,
			"bool-b",
			clip.FlagHelpDefault(""),
		),
		clip.BoolFlag(
			&boolTrue,
			"bool-c",
			clip.FlagHelpDefault("some-default"),
		),
		clip.StringFlag(
			&name,
			"name",
			clip.FlagDescription("Who to greet"),
			clip.FlagEnv("HELLO_NAME"),
		),
		clip.CommandAction(func(*clip.Context) error {
			greeting := fmt.Sprintf("Hello, %s!", name)
			if loud {
				greeting = strings.ToUpper(greeting)
			}
			fmt.Println(greeting)

			return nil
		}),
	)

	args := []string{hello.Name(), "--help"}
	g.NoError(hello.Execute(args))

	g.Should(be.Equal(
		buf.String(),
		`hello - Greet the world

This is a command that will say hello to
a person, or to the world.

Options:
      --bool-a=<bool>
          Default: true

      --bool-b=<bool>

      --bool-c=<bool>
          Default: some-default

  -h, --help
          Print help and exit

  -l, --loud
          Whether to pump up the volume to max.
          This is a very long multi-line, very multi-line
          description of things to come and more!

          Deprecated: Please don't be loud.
          Env: HELLO_LOUD, LOUD

      --name <string>
          Who to greet

          Default: world
          Env: HELLO_NAME
`,
	))
}
