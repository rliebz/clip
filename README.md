# Clip

Clip is a highly opinionated, highly unstable library for building command-line
applications, using [functional options for friendly APIs][functional].

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
  if err := app.Run(os.Args); err != nil {
    log.Fatal(err)
  }
}
```

Of course, this doesn't actually do anything interesting. To create a command:

```go
hello := clip.NewCommand(
  "hello",
  clip.WithDescription("Greet the world"),
  clip.WithBehavior(func(cmd *clip.Command) error {
    fmt.Println("Hello, world!")
    return nil
  }),
)

app := clip.NewCommand(
  "my-app",
  clip.WithCommand(hello),
)

if err := app.Run(os.Args); err != nil {
  log.Fatal(err)
}
```

Arguments and flags can be used as well:

```go
hello := clip.NewCommand(
  "hello",
  clip.WithDescription("Greet a friend"),
  clip.WithFlag(
    clip.NewBoolFlag(
      "loud",
      clip.WithFlagDescription("Whether to pump up the volume to max"),
    ),
  ),
  clip.WithArg(
    clip.NewArg(
      "name",
      clip.WithArgDescription("The person to greet"),
      clip.WithArgValues([]string{"Alice", "Bruce", "Carl"}),
    ),
  ),
  clip.WithBehavior(func(cmd *clip.Command) error {
    greeting := fmt.Sprintf("Hello, %s\n", cmd.Args[0])
    if cmd.Flags["loud"] {
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

if err := app.Run(os.Args); err != nil {
  log.Fatal(err)
}
```


[functional]: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
