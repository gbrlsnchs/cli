package cli

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/tabwriter"
)

var errUnknown = errors.New("unknown command or flag")

// CLI is a command-line interface wrapper that provides flags,
// positional arguments and subcommands.
//
// For providing flags, it uses the flag package from the standard library.
// It overrides some features of the flag package in order to provide a prettier
// help message and to print help to the correct output depending on the situation,
// which is:
//
//   If the user explicitly ask for help by using either -h or -help, it prints help to stdout.
//
//   If the user makes a mistake by missing either a subcommand or a positional argument, it prints
//   help to stderr.
type CLI struct {
	name           string
	entry          *Command
	stdout, stderr io.Writer
	helptxt        string
	codes          struct{ err, misuse int }
}

// New instantiates a new command-line interface with sane defaults,
// which have outputs set to os.Stdout and os.Stderr.
func New(entry *Command, opts ...func(*CLI)) *CLI {
	cli := &CLI{
		entry:   entry,
		stdout:  os.Stdout,
		stderr:  os.Stderr,
		helptxt: "Print this help message.",
		codes:   struct{ err, misuse int }{1, 2},
	}
	for _, o := range opts {
		o(cli)
	}
	return cli
}

// ParseAndRun parses arguments and runs a command, handling flags, subcommands
// and positional arguments according to configuration.
//
// It returns a status code and correctly prints error messages followed by
// usage instructions when necessary.
func (cli *CLI) ParseAndRun(args []string) int {
	return cli.ParseAndRunContext(context.Background(), args)
}

// ParseAndRunContext is just like ParseAndRun but accepts a custom context.
func (cli *CLI) ParseAndRunContext(ctx context.Context, args []string) int {
	if cli.entry == nil {
		panic(fmt.Errorf("%s: cli: nil entry command", cli.name))
	}
	if cli.name == "" {
		// Strip parent directories from the executable's name.
		cli.name = filepath.Base(args[0])
	}
	// This buffer allows printing usage errors with the CLI's name as prefix.
	// Declaring it here prevents from declaring it in every subcommand iteration.
	buf := bytes.NewBufferString(fmt.Sprintf("%s: ", cli.name))
	var (
		code = 0 // success should always be 0, of course
		err  error
		lw   = &lazyWriter{stdout: cli.stdout, stderr: cli.stderr}
	)
	select {
	case <-ctx.Done():
		err = ctx.Err()
	default:
		cli.stdout = (*stdoutWriter)(lw)
		cli.stderr = (*stderrWriter)(lw)
		run := cli.parse(ctx, cli.name, cli.entry, args[1:], buf)
		if run == nil {
			err = errUnknown
		} else {
			err = run()
		}
	}
	// TODO: use custom status codes
	if err != nil {
		if errors.Is(err, errUnknown) {
			lw.flush()
			return cli.codes.misuse
		}
		fmt.Fprintf(lw.stderr, "%v: %v\n", cli.name, err)
		code = cli.codes.err
	}
	lw.flush()
	return code
}

// Stdout is a functional option for creating a CLI that sets w as stdout.
func Stdout(w io.Writer) func(*CLI) {
	return func(cli *CLI) {
		cli.stdout = w
	}
}

// Stdout is a functional option for creating a CLI that sets w as stderr.
func Stderr(w io.Writer) func(*CLI) {
	return func(cli *CLI) {
		cli.stderr = w
	}
}

// HelpDescription changes the default help description message of the CLI.
func HelpDescription(s string) func(*CLI) {
	return func(cli *CLI) {
		cli.helptxt = s
	}
}

// ErrorCode sets a different error code for a CLI. The default is 1.
func ErrorCode(c int) func(*CLI) {
	return func(cli *CLI) {
		cli.codes.err = c
	}
}

// MisuseCode sets a different misuse code for a CLI. The default is 2.
func MisuseCode(c int) func(*CLI) {
	return func(cli *CLI) {
		cli.codes.misuse = c
	}
}

// Name sets a fixed name for the program.
// The default is the first string from parsed args.
func Name(s string) func(*CLI) {
	return func(cli *CLI) {
		cli.name = s
	}
}

func (cli *CLI) parse(ctx context.Context, name string, c *Command, args []string, flagOut *bytes.Buffer) func() error {
	f := flag.NewFlagSet(name, flag.ContinueOnError)
	// Suppress default help messages, since they are printed to stderr even when explicitly requested.
	// See more at https://www.jstorimer.com/blogs/workingwithcode/7766119-when-to-use-stderr-instead-of-stdout.
	var help bool
	if c.Options == nil {
		c.Options = make(map[string]Option, 1)
	}
	c.Options["help"] = BoolOption{
		OptionDetails: OptionDetails{
			Description: cli.helptxt,
			Short:       'h',
		},
		Recipient: &help,
	}
	// Define flags and their aliases to the respective flag set.
	for name, fg := range c.Options {
		fg.Define(f, name)
	}
	// The usage function shows the short, less complete description, in order to not be confuse
	// when a user types a wrong flag.
	f.Usage = func() {
		w := cli.stdout
		if !help {
			w = f.Output()
		}
		tw := tabwriter.NewWriter(w, 0, 0, 4, ' ', 0)
		c.writeUsage(tw, name, help)
		if err := tw.Flush(); err != nil {
			panic(err)
		}
	}
	f.SetOutput(flagOut)
	if err := f.Parse(args); err != nil {
		flagOut.WriteTo(cli.stderr)
		return nil
	}
	if help {
		return usageFunc(f)
	}
	sub := f.Arg(0)
	args = f.Args()
	// Prevent hitting subcommands map when not needed.
	// Also, when there are no subcommands, process args.
	if sub == "" || len(c.Subcommands) == 0 {
		goto exec
	}
	if c, ok := c.Subcommands[sub]; ok {
		return cli.parse(ctx, sub, c, args[1:], flagOut)
	}
	if sub != "" {
		// Bad subcommand.
		cli.printErr(f, fmt.Errorf("command provided but not defined: %s", sub))
		return nil
	}
exec:
	if c.Exec == nil {
		// Command exists BUT has no function attributed to it.
		// This means it should print help to stdout, like Git does.
		help = true // XXX
		return usageFunc(f)
	}
	arglist := new(ArgList)
	if arg := c.Arg; arg != nil {
		arg.AppendTo(arglist)
		if a := arglist.missing(args); a != nil {
			// Bad arguments.
			cli.printErr(f, fmt.Errorf("missing required argument: %s", a.name))
			return nil
		}
		if err := arglist.parse(args); err != nil {
			cli.printErr(f, fmt.Errorf("bad argument parsing: %w", err))
			return nil
		}
	}
	return func() error {
		prg := (*cliMeta)(cli)
		return c.Exec(prg)
	}
}

func (cli *CLI) printErr(f *flag.FlagSet, err error) {
	fmt.Fprintf(cli.stderr, "%s: %v\n\n", cli.name, err)
	f.SetOutput(cli.stderr)
	f.Usage()
}

type cliMeta CLI

func (cli *cliMeta) Name() string      { return cli.name }
func (cli *cliMeta) Stdout() io.Writer { return cli.stdout }
func (cli *cliMeta) Stderr() io.Writer { return cli.stderr }

func usageFunc(f *flag.FlagSet) func() error {
	return func() error {
		f.Usage()
		return nil
	}
}
