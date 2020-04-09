package clitest

import (
	"io"

	"github.com/gbrlsnchs/cli"
)

var _ cli.Program = Program{}
var _ cli.Program = new(Program)

// Program is a stub program that implements cli.Program.
type Program struct {
	name string
	outw io.Writer
	errw io.Writer
}

// NewProgram returns a new stub program.
func NewProgram(name string, outw, errw io.Writer) Program {
	return Program{name, outw, errw}
}

// Name returns the program's name.
func (p Program) Name() string { return p.name }

// Stdout returns the program's stdout.
func (p Program) Stdout() io.Writer { return p.outw }

// Stderr returns the program's stderr.
func (p Program) Stderr() io.Writer { return p.errw }
