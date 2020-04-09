package cli

import (
	"fmt"
	"io"
	"strings"
)

// ExecFunc is a function that receives a program information
// and may return an error, which will be printed to stderr.
type ExecFunc func(Program) error

// Command is a command line command.
//
// If a command doesn't have an Exec function, it is treated as a help command,
// which prints help to stdout.
//
// When a command has one or more subcommands, its Arg will be totally ignored.
type Command struct {
	Description string              // Description describes what the command does.
	Usage       string              // Usage shows how to use the command.
	Exec        ExecFunc            // Exec is the function run by the command.
	Options     map[string]Option   // Options are the command's options (also known as flags).
	Subcommands map[string]*Command // Subcommands store the command's subcommands.
	Arg         Arg                 // Arg is a positional argument.
}

func (c *Command) setDefaultUsage(name string) {
	if c.Usage != "" {
		return
	}
	bd := new(strings.Builder)
	fmt.Fprint(bd, name)
	cmdc := len(c.Subcommands)
	fmt.Fprint(bd, " [<OPTIONS>]")
	if cmdc > 0 {
		fmt.Fprint(bd, " [<COMMAND>]")
	} else if c.Arg != nil {
		c.Arg.WriteDoc(bd)
	}
	c.Usage = bd.String()
}

// Program carries information about a running command.
type Program interface {
	Name() string
	Stdout() io.Writer
	Stderr() io.Writer
}
