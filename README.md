# Clip

Clip is a highly opinionated, highly unstable library for building command-line
applications, using [functional options for friendly APIs][functional].

**Warning: Clip is incomplete software and should not be used by anyone.**

## Quick Start

To get an app running requires no special configuration and does not prescribe
any directory structure. For an app named `my-app`:

```go
package main

import (
  "log"
  "os"

  "github.com/rliebz/clip/cliparg"
  "github.com/rliebz/clip/clipflag"
  "github.com/rliebz/clip/command"
)

func main() {
  // Create a command-line application
  app := command.New("my-app")

  // Run it
  os.Exit(app.Run())
}
```

By default, commands with no action specified print the help documentation:

```text
$ my-app
my-app

Options:
  -h, --help  Print help and exit
```

Of course, since our app doesn't do anything, the help documentation isn't very
useful. Commands can be configured by passing a list of functional options,
such as `command.WithSummary` for a one-line summary, or `command.WithDescription`
for a slightly longer description:

```go
app := command.New(
  "my-app",
  command.WithSummary("A command-line application"),
  command.WithDescription(`This is a simple "Hello World" demo application.`),
)

os.Exit(app.Run())
```

Now when we run it:

```text
$ my-app
my-app - A command-line application

This is a simple "Hello World" demo application.

Options:
  -h, --help  Print help and exit
```

Let's add a sub-command using `WithCommand`, and functionality using
`WithAction`. Because commands are immutable once created, we must declare sub-
commands before their parent commands:

```go
// Define a sub-command "hello"
hello := command.New(
  "hello",
  command.WithSummary("Greet the world"),
  command.WithAction(func(ctx *command.Context) error {
    fmt.Println("Hello, world!")
    return nil
  }),
)

// Create a root command "my-app"
app := command.New(
  "my-app",
  command.WithSummary("A command-line application"),
  command.WithDescription(`This is a simple "Hello World" demo application.`),
  command.WithCommand(hello),
)

// Run it
os.Exit(app.Run())
```

We can see sub-commands in the help documentation:

```text
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

```text
$ my-app hello
Hello, world!
```

### Arguments

By default, arguments are available as a slice from an action's context:

```go
hello := command.New(
  "hello",
  command.WithAction(func (ctx *command.Context) error {
    fmt.Println("Args: ", ctx.Args())
    return nil
  }),
)
```

Generally, however, it is better to explicitly define the arguments. This gives
the benefit of documentation, validation, and tab-completion and can be done using
`command.WithArg` and the `cliparg` package:

```go
name := "World"
hello := command.New(
  "hello",
  command.WithSummary("Greet the world"),
  command.WithArg(
    cliparg.New(
      &name,
      "name",
      cliparg.AsOptional,
      cliparg.WithSummary("The person to greet"),
      cliparg.WithValues([]string{"Alice", "Bruce", "Carl"}),
    ),
  ),
  command.WithAction(func(ctx *command.Context) error {
    greeting := fmt.Sprintf("Hello, %s\n", name)
    fmt.Println(greeting)
    return nil
  }),
)
```

```text
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

### Flags

The `-h`/`--help` flag is defined by default, but more can be created using
`command.WithFlag`/`command.WithActionFlag` and the `clipflag` package.

To create a flag that prints the version and exits, use an action flag:

```go
version := "v0.1.0"
app := command.New(
  "app",
  command.WithActionFlag(
    clipflag.NewToggle(
      "version",
      clipflag.WithShort("V"),
      clipflag.WithSummary("Print the version and exit"),
    ),
    func(ctx *command.Context) error {
      fmt.Println(version)
      return nil
    },
  ),
)
```

```go
loud := false
hello := command.New(
  "hello",
  command.WithFlag(
    clipflag.NewBool(
      &loud,
      "loud",
      clipflag.WithSummary("Whether to pump up the volume to max"),
    ),
  ),
  command.WithAction(func(ctx *command.Context) error {
    if loud {
      fmt.Println("HELLO, WORLD!!!")
    } else {
      fmt.Println("Hello, World!")
    }
    return nil
  }),
)
```

Flags are defined using [POSIX/GNU-style flags][gnu-flags], with `--foo` for a
flag named `"foo"`, and a short, one character flag prefixed with `-` if
specified.

```text
$ my-app hello --help
my-app hello - Greet the world

Usage:
  my-app hello [options]

Options:
  -h, --help  Print help and exit
  -l, --loud  Whether the pump up the volume to max

$ my-app hello
Hello, World!

$ my-app hello --loud
HELLO, WORLD!!!

$ my-app hello -l
HELLO, WORLD!!!
```


[functional]: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
[gnu-flags]: https://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html
