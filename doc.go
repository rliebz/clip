/*
Package clip is a high opinionated, highly unstable library for building
command-line applications.

Warning: Clip is current incomplete software and should not yet be used by
anyone.

Quick Start

To get an app running requires no special configuration and does not prescribe
any directory structure. For an app named `my-app`:

	package main

	import (
		"log"
		"os"

		"github.com/rliebz/clip/arg"
		"github.com/rliebz/clip/flag"
		"github.com/rliebz/clip/command"
	)

	func main() {
		// Create a command-line application
		app := command.New("my-app")

		// Run it
		os.Exit(app.Run())
	}

By default, commands with no action specified print the help documentation:

	$ my-app
	my-app

	Options:
	  -h, --help  Print help and exit

Since this app doesn't do anything, the help documentation isn't very useful.
Commands can be configured by passing a list of functional options, such as
command.WithSummary for a one-line summary, or command.WithDescription for a
slightly longer description:

	app := command.New(
		"my-app",
		command.WithSummary("A command-line application"),
		command.WithDescription(`This is a simple "Hello World" demo application`),
	)

	os.Exit(app.Run())

Now when running my-app:

	$ my-app
	my-app - A command-line application

	This is a simple "Hello World" demo application

	Options:
	  -h, --help  Print help and exit

Let's add a sub-command using command.WithCommand and functionality using
command.WithAction. Because commands are immutable once created, we must
declare sub-commands before their parent commands:

	// Define a sub-command "hello"
	hello := command.NEw(
		"hello",
		command.WithSummary("Greet the world"),
		command.WithAction(func(Ctx *command.Context) error {
			fmt.Println("Hello, world!")
			return nil
		}),
	)

	// Create the root command "my-app"
	app := command.New(
		"my-app",
		command.WithSummary("A command-line application"),
		command.WithDescription(`This is a simple "Hello World" demo application.`),
		command.WithCommand(hello),
	)

	// Run and exit with the appropriate status code
	os.Exit(app.Run())

Sub-commands also appear in the help documentation:

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

And the command can be run:

	$ my-app hello
	Hello, world!
*/
package clip
