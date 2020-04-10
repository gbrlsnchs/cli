package cli

import (
	"fmt"
	"io"
	"sort"
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
	Exec        ExecFunc            // Exec is the function run by the command.
	Options     map[string]Option   // Options are the command's options (also known as flags).
	Subcommands map[string]*Command // Subcommands store the command's subcommands.
	Arg         Arg                 // Arg is a positional argument.
}

func (c *Command) writeUsage(w io.Writer, name string, showDesc bool) {
	// DESCRIPTION
	if showDesc && c.Description != "" {
		wrapWrite(w, c.Description)
		fmt.Fprint(w, "\n\n")
	}
	// USAGE (A.K.A. SUMMARY)
	fmt.Fprintln(w, "USAGE:")
	fmt.Fprintf(w, "\t%s [OPTIONS]", name)
	nsub := len(c.Subcommands)
	if nsub > 0 {
		fmt.Fprint(w, " ")
		cstart, cend := "<", ">"
		if c.Exec != nil {
			cstart, cend = "[", "]"
		}
		fmt.Fprintf(w, "%sCOMMAND%s", cstart, cend)
	} else if arg := c.Arg; arg != nil {
		arg.WriteDoc(w)
	}
	fmt.Fprint(w, "\n\nOPTIONS:\n") // this is always printed, since help option is always present
	// OPTIONS
	copts := c.Options
	optl := make([]string, 0, len(copts))
	for name := range copts {
		optl = append(optl, name)
	}
	sort.Strings(optl)
	for _, o := range optl {
		fmt.Fprint(w, "\t")
		copts[o].WriteDoc(w, o)
		fmt.Fprintln(w)
	}
	// COMMANDS
	if nsub > 0 {
		fmt.Fprint(w, "\nCOMMANDS:\n")
		ccmds := c.Subcommands
		subl := make([]string, 0, nsub)
		for name := range ccmds {
			subl = append(subl, name)
		}
		sort.Strings(subl)
		for _, c := range subl {
			fmt.Fprintf(w, "\t%s", c)
			if desc := ccmds[c].Description; desc != "" {
				fmt.Fprintf(w, "\t%s", desc)
			}
			fmt.Fprintln(w)
		}
	}
}

// Program carries information about a running command.
type Program interface {
	Name() string
	Stdout() io.Writer
	Stderr() io.Writer
}

func wrapWrite(w io.Writer, line string) {
	paragraphs := strings.Split(line, "\n")
	for i, p := range paragraphs {
		words := strings.Split(p, " ")
		first := words[0]
		n := len(first)
		const limit = 72
		fmt.Fprint(w, words[0])
		for i, s := range words[1:] {
			spacing := " "
			n += 1 + len(s) // spacing + word
			if n > limit {
				n = 0 // reset characters count
				spacing = "\n"
			} else if i == len(words)-1 {
				spacing = ""
			}
			fmt.Fprintf(w, "%s%s", spacing, s)
		}
		if i == len(paragraphs)-1 {
			continue
		}
		fmt.Fprintln(w)
	}
}
