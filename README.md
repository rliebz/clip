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
  app := clip.NewCommand("my-app")
  os.Exit(app.Run())
}
```

By default, commands with no action specified print the help documentation:

```text
$ my-app
my-app
```

Of course, this doesn't actually do anything interesting. To create a command:

```go
hello := clip.NewCommand(
  "hello",
  clip.WithSummary("Greet the world"),
  clip.WithAction(func(ctx *clip.Context) error {
    fmt.Println("Hello, world!")
    return nil
  }),
)

app := clip.NewCommand(
  "my-app",
  clip.WithSummary("A demo application"),
  clip.WithCommand(hello),
)

os.Exit(app.Run())
```

Now that there is functionality, the help documentation is more useful:

```text
$ my-app
my-app - A demo application

Commands:
  hello  Greet the world
```

And the command can be run:

```text
$ my-app hello
Hello, world!
```

Arguments and flags can be used as well:

```go
loud := false
hello := clip.NewCommand(
  "hello",
  clip.WithSummary("Greet a friend"),
  clip.WithFlag(
    clipflag.NewBool(
      &loud,
      "loud",
      clipflag.WithSummary("Whether to pump up the volume to max"),
    ),
  ),
  clip.WithArg(
    cliparg.New(
      "name",
      cliparg.WithSummary("The person to greet"),
      cliparg.WithValues([]string{"Alice", "Bruce", "Carl"}),
    ),
  ),
  clip.WithAction(func(ctx *clip.Context) error {
    greeting := fmt.Sprintf("Hello, %s\n", ctx.Args()[0])
    if loud {
      greeting = strings.ToUpper(greeting)
    }
    fmt.Println(greeting)
    return nil
  }),
)

app := clip.NewCommand(
  "my-app",
  clip.WithCommand(hello),
)

os.Exit(app.Run())
```


[functional]: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
