# Clip

Clip is a highly opinionated, highly unstable library for building command-line
applications.

**Warning: Clip is currently incomplete software and should not yet be used by
anyone.**

## Quick Start

To get an app running requires no special configuration and does not prescribe
any directory structure. For an app named `my-app`:

```go
package main

import (
	"log"
	"os"

	"github.com/rliebz/clip"
)

func main() {
	// Create a command-line application
	app := clip.NewCommand("my-app")

	// Run it
	os.Exit(app.Run())
}
```

By default, commands with no action specified print the help documentation:

```console
$ my-app
my-app

Options:
  -h, --help  Print help and exit
```

Since this app doesn't do anything, the help documentation isn't very useful.
Commands can be configured by passing a list of functional options, such as
`clip.CommandSummary` for a one-line summary, or `clip.CommandDescription` for
a slightly longer description:

```go
app := clip.NewCommand(
	"my-app",
	clip.CommandSummary("A command-line application"),
	clip.CommandDescription(`This is a simple "Hello World" demo application.`),
)

os.Exit(app.Run())
```

Now when running my-app:

```console
$ my-app
my-app - A command-line application

This is a simple "Hello World" demo application.

Options:
  -h, --help  Print help and exit
```

Let's add a sub-command using `clip.CommandSubCommand` and functionality using
`clip.CommandAction`. Because commands are immutable once created, we must
declare sub-commands before their parent commands:

```go
// Define a sub-command "hello"
hello := clip.NewCommand(
	"hello",
	clip.CommandSummary("Greet the world"),
	clip.CommandAction(func(ctx *clip.Context) error {
	  fmt.Println("Hello, world!")
	  return nil
	}),
)

// Create the root command "my-app"
app := clip.NewCommand(
	"my-app",
	clip.CommandSummary("A command-line application"),
	clip.CommandDescription(`This is a simple "Hello World" demo application.`),
	clip.CommandSubCommand(hello),
)

// Run it
os.Exit(app.Run())
```

Sub-commands also appear in the help documentation:

```console
$ my-app
my-app - A command-line application

This is a simple "Hello World" demo application.

Commands:
  hello  Greet the world

Options:
  -h, --help  Print help and exit

$ my-app hello --help
my-app hello - Greet the world

Options:
  -h, --help  Print help and exit
```

And the command can be run:

```console
$ my-app hello
Hello, world!
```

### Arguments

By default, any unexpected arguments passed to a command are considered an
error.

To make arguments available as a slice from an action's context, the function
`clip.CommandArgs` can be used:

```go
var args []string
hello := clip.NewCommand(
	"hello",
	clip.CommandArgs(&args),
	clip.CommandAction(func (ctx *clip.Context) error {
	  fmt.Println("Args: ", args)
	  return nil
	}),
)
```

Generally, however, it is better to explicitly define the arguments. This gives
the benefit of documentation, validation, and tab-completion.

```go
name := "World"
hello := clip.NewCommand(
	"hello",
	clip.CommandSummary("Greet the world"),
	clip.StringArg(
    &name,
    "name",
    clip.ArgOptional,
    clip.ArgSummary("The person to greet"),
    clip.ArgValues([]string{"Alice", "Bruce", "Carl"}),
	),
	clip.CommandAction(func(ctx *clip.Context) error {
		greeting := fmt.Sprintf("Hello, %s\n", name)
		fmt.Println(greeting)
		return nil
	}),
)
```

This produces an app with the following behavior:

```console
$ my-app hello --help
my-app hello - Greet the world

Usage:
  my-app hello [<name>]

Args:
  name  The person to greet
        One of: Alice, Bruce, Carl

Options:
  -h, --help  Print help and exit

$ my-app hello
Hello, World!

$ my-app hello Alice
Hello, Alice!

$ my-app hello Alex
Error: argument "Alex" must be one of: Alice, Bruce, Carl
```

Arguments and sub-commands are mutually exclusive.

### Flags

The `-h`/`--help` flag is defined by default.

To create a flag that prints the version and exits, use an action flag:

```go
version := "v0.1.0"
app := clip.NewCommand(
	"app",
	clip.ToggleFlag(
		"version",
		clip.FlagShort("V"),
		clip.FlagDescription("Print the version and exit"),
		clip.FlagAction(func(ctx *clip.Context) error {
			fmt.Println(version)
			return nil
		})
	),
)
```

Flags can be created with different types, such as bool and string:

```go
loud := false
name := "world"
hello := clip.NewCommand(
	"hello",
	clip.CommandSummary("Greet the world"),
	clip.BoolFlag(
		&loud,
		"loud",
		clip.FlagDescription("Whether to pump up the volume to max"),
	),
	clip.StringFlag(
		&name,
		"name",
		clip.FlagShort("l"),
		clip.FlagDescription("Who to greet"),
	),
	clip.CommandAction(func(ctx *clip.Context) error {
		greeting := fmt.Sprintf("Hello, %s!", name)
		if loud {
			greeting = strings.ToUpper(greeting)
		}
		fmt.Println(greeting)

		return nil
	}),
)
```

Flags are defined using [POSIX/GNU-style flags][gnu-flags], with `--foo` for a
flag named `"foo"`, and a short, one character flag prefixed with `-` if
specified.

```console
$ my-app hello --help
my-app hello - Greet the world

Usage:
  my-app hello [options]

Flags:
  -l, --loud  Whether to pump up the volume to max
      --name  Who to greet
  -h, --help  Print help and exit

$ my-app hello
Hello, World!

$ my-app hello --loud
HELLO, WORLD!!!

$ my-app hello -l
HELLO, WORLD!!!
```


[gnu-flags]: https://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html
