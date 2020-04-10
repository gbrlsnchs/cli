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

// copy returns a copy of cfg, thus not allowing the original configuration
// to be modified by registering functions.
// This is one way to simulate a const reference, which, although modifiable,
// won't interfere with the original instance, so modifying it is harmless.
func (cfg *appConfig) copy() appConfig { return *cfg }

// helloCmd says hello to someone. Toggle upper to scream.
type helloCmd struct {
	name  string
	upper bool
}

func (cmd *helloCmd) register(getcfg func() appConfig) cli.ExecFunc {
	return func(prg cli.Program) error {
		appcfg := getcfg()
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

func (cmd *concatCmd) register(getcfg func() appConfig) cli.ExecFunc {
	return func(prg cli.Program) error {
		appcfg := getcfg()
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

func (cmd *joinCmd) register(getcfg func() appConfig) cli.ExecFunc {
	return func(prg cli.Program) error {
		appcfg := getcfg()
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
		appcfg appConfig
	)
	cmdl := cli.New(&cli.Command{
		Description: `This is a simple program that serves as an example for how to use package cli.

Its commands should not be taken seriously, since they do nothing really great, but serve well for demonstration.`,
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
				Exec: root.hello.register(appcfg.copy),
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
				Exec: root.concat.register(appcfg.copy),
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
				Exec: root.join.register(appcfg.copy),
			},
		},
	}, cli.Name("my-cmd"))
	code := cmdl.ParseAndRun(os.Args)
	os.Exit(code)
}
