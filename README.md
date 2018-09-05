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

```text
$ my-app
my-app

Options:
  -h, --help  Print help and exit
```

Of course, since our app doesn't do anything, the help documentation isn't very
useful. Commands can be configured by passing a list of functional options,
such as `clip.WithSummary` for a one-line summary, or `clip.WithDescription`
for a slightly longer description:

```go
app := clip.NewCommand(
  "my-app",
  clip.WithSummary("A command-line application"),
  clip.WithDescription(`This is a simple "Hello World" demo application.`),
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
hello := clip.NewCommand(
  "hello",
  clip.WithSummary("Greet the world"),
  clip.WithAction(func(ctx *clip.Context) error {
    fmt.Println("Hello, world!")
    return nil
  }),
)

// Create a root command "my-app"
app := clip.NewCommand(
  "my-app",
  clip.WithSummary("A command-line application"),
  clip.WithDescription(`This is a simple "Hello World" demo application.`),
  clip.WithCommand(hello),
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
hello := clip.NewCommand(
  "hello",
  clip.WithAction(func (ctx *Context) error {
    fmt.Println("Args: ", ctx.Args())
    return nil
  }),
)
```

Generally, however, it is better to explicitly define the arguments. This gives
the benefit of documentation, validation, and tab-completion and can be done using
`clip.WithArg` and the `cliparg` package:

```go
name := "World"
hello := clip.NewCommand(
  "hello",
  clip.WithSummary("Greet the world"),
  clip.WithArg(
    cliparg.New(
      &name,
      "name",
      cliparg.AsOptional,
      cliparg.WithSummary("The person to greet"),
      cliparg.WithValues([]string{"Alice", "Bruce", "Carl"}),
    ),
  ),
  clip.WithAction(func(ctx *clip.Context) error {
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
`clip.WithFlag` and the `clipflag` package:

```go
loud := false
hello := clip.NewCommand(
  "hello",
  clip.WithSummary("Greet the world"),
  clip.WithFlag(
    clipflag.NewBool(
      &loud,
      "loud",
      clipflag.WithShort("l"),
      clipflag.WithSummary("Whether to pump up the volume to max"),
    ),
  ),
  clip.WithAction(func(ctx *clip.Context) error {
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
