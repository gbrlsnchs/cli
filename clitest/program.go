package clitest

import (
	"io"
	"strings"

	"github.com/gbrlsnchs/cli"
)

var _ cli.Program = Program{}
var _ cli.Program = new(Program)

// Program is a stub program that implements cli.Program.
type Program struct {
	name   string
	comb   *strings.Builder
	out    *strings.Builder
	errout *strings.Builder
}

// NewProgram returns a new stub program.
func NewProgram(name string) Program {
	comb := new(strings.Builder)
	out := new(strings.Builder)
	errw := new(strings.Builder)
	return Program{name, comb, out, errw}
}

// Name returns the program's name.
func (p Program) Name() string { return p.name }

// Stdout returns the program's stdout.
func (p Program) Stdout() io.Writer { return io.MultiWriter(p.out, p.comb) }

// Stderr returns the program's stderr.
func (p Program) Stderr() io.Writer { return io.MultiWriter(p.errout, p.comb) }

// Output returns what has been written to standard output.
func (p Program) Output() string { return p.out.String() }

// ErrOutput returns what has been written to standard error.
func (p Program) ErrOutput() string { return p.errout.String() }

// CombinedOutput returns what has been written to both standard output and error.
func (p Program) CombinedOutput() string { return p.comb.String() }
