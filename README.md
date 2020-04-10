# cli
![Linux, macOS and Windows](https://github.com/gbrlsnchs/cli/workflows/Linux,%20macOS%20and%20Windows/badge.svg)
[![GoDoc](https://godoc.org/github.com/gbrlsnchs/cli?status.svg)](https://godoc.org/github.com/gbrlsnchs/cli)

## About
This is a library that adds some CLI functionalities on top of [Go's flag package] while preserving the single dash Go-style flags. Some of those functionalities are:
- Subcommands
- Positional arguments
  - Both required and optional arguments
  - Repeating arguments
- More robust help message
- Correct handling of help flags
  - Print help to stdout when help is explicitly requested (via `-h` or `-help` options)
  - Print help to stderr when the CLI is misused (by requesting a bad command or argument)
- Prefix errors with the program's name
- Easy to set up (the whole CLI can be configured all at once)

### Principles
This library enforces some principles that might not suit everybody's taste:
- Go-style flags (single dash for both long and short options)
- No flags with short name only
- No conditional flags (internal parser is still [Go's flag package])
- No required flags (if it's required, make it an argument)
- Flags must come right after its command, before args
- Do not expose types from external packages
- Boring (it's basically configuration plus your logic)

Really, if you're looking for GNU-style flags with the whole parsing kung fu, this library might not be for you. But don't worry, Go is full of awesome libraries out there for such purposes.

## Example
### Multiple subcommands with top-level configuration
```go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gbrlsnchs/cli"
)

// appConfig is the program's configuration.
// It can be used by any command.
type appConfig struct {
	quiet bool
}

// helloCmd says hello to someone. Toggle upper to scream.
type helloCmd struct {
	name  string
	upper bool
}

func (cmd *helloCmd) register(appcfg *appConfig) cli.ExecFunc {
	return func(prg cli.Program) error {
		s := cmd.name
		if cmd.upper {
			s = strings.ToUpper(s)
		}
		if !appcfg.quiet {
			fmt.Fprintf(prg.Stdout(), "Hello, %s!\n", s)
		}
		return nil
	}
}

// concatCmd prints the concatenation of two words.
type concatCmd struct {
	first  string
	second string
}

func (cmd *concatCmd) register(appcfg *appConfig) cli.ExecFunc {
	return func(prg cli.Program) error {
		if !appcfg.quiet {
			fmt.Fprintf(prg.Stdout(), "%s %s\n", cmd.first, cmd.second)
		}
		return nil
	}
}

// joinCmd prints a list of words joined by a separator.
type joinCmd struct {
	words []string
	sep   string
}

func (cmd *joinCmd) register(appcfg *appConfig) cli.ExecFunc {
	return func(prg cli.Program) error {
		if !appcfg.quiet {
			s := strings.Join(cmd.words, cmd.sep)
			fmt.Fprintln(prg.Stdout(), s)
		}
		return nil
	}
}

// rootCmd is simply a store for commands, and also a way to
// initialize all of them at once by using rootCmd's zero value.
type rootCmd struct {
	hello  helloCmd
	concat concatCmd
	join   joinCmd
}

func main() {
	var (
		root   rootCmd
		appcfg = new(appConfig)
	)
	cmdl := cli.New(&cli.Command{
		Description: `This is a simple program that serves as an example for how to use package cli.

Its commands should not be taken seriously, since they do nothing really great, but serve well for demonstration`,
		Options: map[string]cli.Option{
			"quiet": cli.BoolOption{
				OptionDetails: cli.OptionDetails{
					Description: "Turn output off.",
					Short:       'q',
				},
				Recipient: &appcfg.quiet,
			},
		},
		Subcommands: map[string]*cli.Command{
			"hello": {
				Description: "Say hello to someone.",
				Arg: cli.StringArg{
					Label:     "NAME",
					Required:  true,
					Recipient: &root.hello.name,
				},
				Options: map[string]cli.Option{
					"upper": cli.BoolOption{
						OptionDetails: cli.OptionDetails{
							Description: "Convert name to uppercase.",
						},
						Recipient: &root.hello.upper,
					},
				},
				Exec: root.hello.register(appcfg),
			},
			"concat": {
				Description: "Concatenate two words.",
				Arg: cli.StringArg{
					Label:     "FIRST WORD",
					Required:  true,
					Recipient: &root.concat.first,
					Next: cli.StringArg{
						Label:     "LAST WORD",
						Required:  true,
						Recipient: &root.concat.second,
					},
				},
				Exec: root.concat.register(appcfg),
			},
			"join": {
				Description: "Join strings together.",
				Arg: cli.RepeatingArg{
					Label:     "WORD",
					Required:  true,
					Recipient: &root.join.words,
				},
				Options: map[string]cli.Option{
					"separator": cli.StringOption{
						OptionDetails: cli.OptionDetails{
							Description: "Set a custom separator.",
							Short:       's',
							ArgLabel:    "SEPARATOR",
						},
						Recipient: &root.join.sep,
						DefValue:  ",",
					},
				},
				Exec: root.join.register(appcfg),
			},
		},
	}, cli.Name("my-cmd"))
	code := cmdl.ParseAndRun(os.Args)
	os.Exit(code)
}
```

#### Outputs
<details><summary>Running the main program (prints help)</summary>
<p>

```shell
$ go run .examples/full_featured/main.go -h
This is a simple program that serves as an example for how to use
package cli.

Its commands should not be taken seriously, since they do nothing really
great, but serve well for demonstration

USAGE:
    my-cmd [OPTIONS] <COMMAND>

OPTIONS:
    -h, -help     Print this help message.
    -q, -quiet    Turn output off.

COMMANDS:
    concat    Concatenate two words.
    hello     Say hello to someone.
    join      Join strings together.
```

</p>
</details>

<details><summary>Help for the "hello" program</summary>
<p>

```shell
$ go run .examples/full_featured/main.go hello -h
Say hello to someone.

USAGE:
    hello [OPTIONS] <NAME>

OPTIONS:
    -h, -help     Print this help message.
        -upper    Convert name to uppercase.
```

</p>
</details>

<details><summary>Help for the "concat" program</summary>
<p>

```shell
$ go run .examples/full_featured/main.go concat -h
Concatenate two words.

USAGE:
    concat [OPTIONS] <FIRST WORD> <LAST WORD>

OPTIONS:
    -h, -help    Print this help message.
```

</p>
</details>

<details><summary>Help for the "join" program</summary>
<p>

```shell
$ go run .examples/full_featured/main.go join -h
Join strings together.

USAGE:
    join [OPTIONS] <WORD> [...]

OPTIONS:
    -h, -help                     Print this help message.
    -s, -separator <SEPARATOR>    Set a custom separator. (default: ",")
```

</p>
</details>

[Go's flag package]: https://golang.org/pkg/flag/
