package cli_test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/gbrlsnchs/cli"
	"github.com/google/go-cmp/cmp"
)

func TestCommandLine(t *testing.T) {
	t.Run("ParseAndRun", testCommandLineParseAndRun)
}

func testCommandLineParseAndRun(t *testing.T) {
	type testCommand struct {
		*testing.T
		parg1 string
		parg2 string
		parg3 string
		rargs []string
		fstr  string
		fbool bool
		fint  int
		flong int64
	}
	var root testCommand
	testCases := []struct {
		desc         string
		entry        *cli.Command
		opts         []func(*cli.CLI)
		args         []string
		wantCode     int
		wantOut      string
		wantErr      string
		wantCombined string
	}{
		{
			desc: "main command without args or options",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "", root.parg1; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					return nil
				},
				Subcommands: nil,
			},
			args:         []string{"test"},
			wantCode:     0,
			wantOut:      "",
			wantErr:      "",
			wantCombined: "",
		},
		{
			desc: "print help of main command without args or options",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "", root.parg1; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					return nil
				},
				Subcommands: nil,
			},
			args:     []string{"test", "-h"},
			wantCode: 0,
			wantOut: `USAGE:
    test [<OPTIONS>]

OPTIONS:
    -h, -help    print help information
`,
			wantErr: "",
			wantCombined: `USAGE:
    test [<OPTIONS>]

OPTIONS:
    -h, -help    print help information
`,
		},
		{
			desc: "main command with one optional arg",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "foo", root.parg1; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					return nil
				},
				Subcommands: nil,
				Arg: cli.StringArg{
					Label:     "FOO",
					Recipient: &root.parg1,
				},
			},
			args:         []string{"test", "foo"},
			wantCode:     0,
			wantOut:      "",
			wantErr:      "",
			wantCombined: "",
		},
		{
			desc: "print help of main command with one optional arg",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "", root.parg1; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					return nil
				},
				Subcommands: nil,
				Arg: cli.StringArg{
					Label:     "FOO",
					Recipient: &root.parg1,
				},
			},
			args:     []string{"test", "-h"},
			wantCode: 0,
			wantOut: `USAGE:
    test [<OPTIONS>] [<FOO>]

OPTIONS:
    -h, -help    print help information
`,
			wantErr: "",
			wantCombined: `USAGE:
    test [<OPTIONS>] [<FOO>]

OPTIONS:
    -h, -help    print help information
`,
		},
		{
			desc: "main command missing one optional arg",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "", root.parg1; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					return nil
				},
				Subcommands: nil,
				Arg: cli.StringArg{
					Label:     "FOO",
					Recipient: &root.parg1,
				},
			},
			args:         []string{"test"},
			wantCode:     0,
			wantOut:      "",
			wantErr:      "",
			wantCombined: "",
		},
		{
			desc: "main command with one required arg",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "foo", root.parg1; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					return nil
				},
				Subcommands: nil,
				Arg: cli.StringArg{
					Label:     "FOO",
					Required:  true,
					Recipient: &root.parg1,
				},
			},
			args:         []string{"test", "foo"},
			wantCode:     0,
			wantOut:      "",
			wantErr:      "",
			wantCombined: "",
		},
		{
			desc: "print help of main command with one required arg",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "", root.parg1; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					return nil
				},
				Subcommands: nil,
				Arg: cli.StringArg{
					Label:     "FOO",
					Required:  true,
					Recipient: &root.parg1,
				},
			},
			args:     []string{"test", "-h"},
			wantCode: 0,
			wantOut: `USAGE:
    test [<OPTIONS>] <FOO>

OPTIONS:
    -h, -help    print help information
`,
			wantErr: "",
			wantCombined: `USAGE:
    test [<OPTIONS>] <FOO>

OPTIONS:
    -h, -help    print help information
`,
		},
		{
			desc: "main command missing one required arg",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "", root.parg1; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					return nil
				},
				Subcommands: nil,
				Arg: cli.StringArg{
					Label:     "FOO",
					Required:  true,
					Recipient: &root.parg1,
				},
			},
			args:     []string{"test"},
			wantCode: 2,
			wantOut:  "",
			wantErr: `test: missing required argument: FOO
USAGE:
    test [<OPTIONS>] <FOO>

OPTIONS:
    -h, -help    print help information
`,
			wantCombined: `test: missing required argument: FOO
USAGE:
    test [<OPTIONS>] <FOO>

OPTIONS:
    -h, -help    print help information
`,
		},
		{
			desc: "main command with two optional args",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "foo", root.parg1; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					if want, got := "bar", root.parg2; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					return nil
				},
				Subcommands: nil,
				Arg: cli.StringArg{
					Label:     "FOO",
					Recipient: &root.parg1,
					Next: cli.StringArg{
						Label:     "BAR",
						Recipient: &root.parg2,
					},
				},
			},
			args:         []string{"test", "foo", "bar"},
			wantCode:     0,
			wantOut:      "",
			wantErr:      "",
			wantCombined: "",
		},
		{
			desc: "print help of main command with two optional args",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "", root.parg1; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					if want, got := "", root.parg2; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					return nil
				},
				Subcommands: nil,
				Arg: cli.StringArg{
					Label:     "FOO",
					Recipient: &root.parg1,
					Next: cli.StringArg{
						Label:     "BAR",
						Recipient: &root.parg2,
					},
				},
			},
			args:     []string{"test", "-h"},
			wantCode: 0,
			wantOut: `USAGE:
    test [<OPTIONS>] [<FOO> [<BAR>]]

OPTIONS:
    -h, -help    print help information
`,
			wantErr: "",
			wantCombined: `USAGE:
    test [<OPTIONS>] [<FOO> [<BAR>]]

OPTIONS:
    -h, -help    print help information
`,
		},
		{
			desc: "main command missing two optional args",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "", root.parg1; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					if want, got := "", root.parg2; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					return nil
				},
				Subcommands: nil,
				Arg: cli.StringArg{
					Label:     "FOO",
					Recipient: &root.parg1,
					Next: cli.StringArg{
						Label:     "BAR",
						Recipient: &root.parg2,
					},
				},
			},
			args:         []string{"test"},
			wantCode:     0,
			wantOut:      "",
			wantErr:      "",
			wantCombined: "",
		},
		{
			desc: "main command with first arg required and second one optional",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "foo", root.parg1; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					if want, got := "bar", root.parg2; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					return nil
				},
				Subcommands: nil,
				Arg: cli.StringArg{
					Label:     "FOO",
					Required:  true,
					Recipient: &root.parg1,
					Next: cli.StringArg{
						Label:     "BAR",
						Recipient: &root.parg2,
					},
				},
			},
			args:         []string{"test", "foo", "bar"},
			wantCode:     0,
			wantOut:      "",
			wantErr:      "",
			wantCombined: "",
		},
		{
			desc: "print help of main command with first arg required and second one optional",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "", root.parg1; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					if want, got := "", root.parg2; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					return nil
				},
				Subcommands: nil,
				Arg: cli.StringArg{
					Label:     "FOO",
					Required:  true,
					Recipient: &root.parg1,
					Next: cli.StringArg{
						Label:     "BAR",
						Recipient: &root.parg2,
					},
				},
			},
			args:     []string{"test", "-h"},
			wantCode: 0,
			wantOut: `USAGE:
    test [<OPTIONS>] <FOO> [<BAR>]

OPTIONS:
    -h, -help    print help information
`,
			wantErr: "",
			wantCombined: `USAGE:
    test [<OPTIONS>] <FOO> [<BAR>]

OPTIONS:
    -h, -help    print help information
`,
		},
		{
			desc: "main command missing second optional argument with first one required",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "foo", root.parg1; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					if want, got := "", root.parg2; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					return nil
				},
				Subcommands: nil,
				Arg: cli.StringArg{
					Label:     "FOO",
					Recipient: &root.parg1,
					Next: cli.StringArg{
						Label:     "BAR",
						Recipient: &root.parg2,
					},
				},
			},
			args:         []string{"test", "foo"},
			wantCode:     0,
			wantOut:      "",
			wantErr:      "",
			wantCombined: "",
		},
		{
			desc: "main command with one optional repeating arg",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := []string{"foo"}, root.rargs; !cmp.Equal(got, want) {
						t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
					}
					return nil
				},
				Subcommands: nil,
				Arg: cli.RepeatingArg{
					Label:     "FOO",
					Recipient: &root.rargs,
				},
			},
			args:         []string{"test", "foo"},
			wantCode:     0,
			wantOut:      "",
			wantErr:      "",
			wantCombined: "",
		},
		{
			desc: "main command with two optional repeating arg",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := []string{"foo", "bar"}, root.rargs; !cmp.Equal(got, want) {
						t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
					}
					return nil
				},
				Subcommands: nil,
				Arg: cli.RepeatingArg{
					Label:     "FOO",
					Recipient: &root.rargs,
				},
			},
			args:         []string{"test", "foo", "bar"},
			wantCode:     0,
			wantOut:      "",
			wantErr:      "",
			wantCombined: "",
		},
		{
			desc: "print help of main command with one optional repeating arg",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := ([]string)(nil), root.rargs; !cmp.Equal(got, want) {
						t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
					}
					return nil
				},
				Subcommands: nil,
				Arg: cli.RepeatingArg{
					Label:     "FOO",
					Recipient: &root.rargs,
				},
			},
			args:     []string{"test", "-h"},
			wantCode: 0,
			wantOut: `USAGE:
    test [<OPTIONS>] [<FOO>...]

OPTIONS:
    -h, -help    print help information
`,
			wantErr: "",
			wantCombined: `USAGE:
    test [<OPTIONS>] [<FOO>...]

OPTIONS:
    -h, -help    print help information
`,
		},
		{
			desc: "main command missing one optional repeating arg",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := ([]string)(nil), root.rargs; !cmp.Equal(got, want) {
						t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
					}
					return nil
				},
				Subcommands: nil,
				Arg: cli.StringArg{
					Label:     "FOO",
					Recipient: &root.parg1,
				},
			},
			args:         []string{"test"},
			wantCode:     0,
			wantOut:      "",
			wantErr:      "",
			wantCombined: "",
		},
		{
			desc: "main command with options",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "foo", root.fstr; got != want {
						t.Errorf("want %q, got %q", want, got)
					}
					if want, got := true, root.fbool; got != want {
						t.Errorf("want %t, got %t", want, got)
					}
					if want, got := 0xDEADC0DE, root.fint; got != want {
						t.Errorf("want %d, got %d", want, got)
					}
					if want, got := int64(0xABADBABE), root.flong; got != want {
						t.Errorf("want %d, got %d", want, got)
					}
					return nil
				},
				Subcommands: nil,
				Options: map[string]cli.Option{
					"string": cli.StringOption{
						OptionDetails: cli.OptionDetails{
							Description: "pass a string here",
							Short:       's',
							ArgLabel:    "TEXT",
						},
						DefValue:  "bar",
						Recipient: &root.fstr,
					},
					"bool": cli.BoolOption{
						OptionDetails: cli.OptionDetails{
							Description: "pass a boolean here",
							Short:       'b',
							ArgLabel:    "TRUE|FALSE",
						},
						Recipient: &root.fbool,
					},
					"int": cli.IntOption{
						OptionDetails: cli.OptionDetails{
							Description: "pass an integer here",
							Short:       'i',
							ArgLabel:    "NUMBER",
						},
						Recipient: &root.fint,
					},
					"int64": cli.Int64Option{
						OptionDetails: cli.OptionDetails{
							Description: "pass a 64-bit integer here",
							Short:       'I',
							ArgLabel:    "64-BIT NUMBER",
						},
						Recipient: &root.flong,
					},
					"dull": cli.StringOption{
						OptionDetails: cli.OptionDetails{
							Description: "just a dull flag",
						},
						Recipient: newDullStr(),
					},
				},
			},
			args:         []string{"test", "-s", "foo", "-b", "-i", "0xDEADC0DE", "-I", "0xABADBABE"},
			wantCode:     0,
			wantOut:      "",
			wantErr:      "",
			wantCombined: "",
		},
		{
			desc: "print help of main command with one optional arg",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "foo", root.fstr; got != want {
						t.Errorf("want %q, got %q", want, got)
					}
					if want, got := true, root.fbool; got != want {
						t.Errorf("want %t, got %t", want, got)
					}
					if want, got := 0xDEADC0DE, root.fint; got != want {
						t.Errorf("want %d, got %d", want, got)
					}
					if want, got := int64(0xABADBABE), root.flong; got != want {
						t.Errorf("want %d, got %d", want, got)
					}
					return nil
				},
				Subcommands: nil,
				Options: map[string]cli.Option{
					"string": cli.StringOption{
						OptionDetails: cli.OptionDetails{
							Description: "pass a string here",
							Short:       's',
							ArgLabel:    "TEXT",
						},
						DefValue:  "bar",
						Recipient: &root.fstr,
					},
					"bool": cli.BoolOption{
						OptionDetails: cli.OptionDetails{
							Description: "pass a boolean here",
							Short:       'b',
							ArgLabel:    "TRUE|FALSE",
						},
						Recipient: &root.fbool,
					},
					"int": cli.IntOption{
						OptionDetails: cli.OptionDetails{
							Description: "pass an integer here",
							Short:       'i',
							ArgLabel:    "NUMBER",
						},
						Recipient: &root.fint,
					},
					"int64": cli.Int64Option{
						OptionDetails: cli.OptionDetails{
							Description: "pass a 64-bit integer here",
							Short:       'I',
							ArgLabel:    "64-BIT NUMBER",
						},
						Recipient: &root.flong,
					},
					"dull": cli.StringOption{
						OptionDetails: cli.OptionDetails{
							Description: "just a dull option",
						},
						Recipient: newDullStr(),
					},
				},
			},
			args:     []string{"test", "-h"},
			wantCode: 0,
			wantOut: `USAGE:
    test [<OPTIONS>]

OPTIONS:
    -b, -bool <TRUE|FALSE>        pass a boolean here
        -dull                     just a dull option
    -h, -help                     print help information
    -i, -int <NUMBER>             pass an integer here
    -I, -int64 <64-BIT NUMBER>    pass a 64-bit integer here
    -s, -string <TEXT>            pass a string here (default: "bar")
`,
			wantErr: "",
			wantCombined: `USAGE:
    test [<OPTIONS>]

OPTIONS:
    -b, -bool <TRUE|FALSE>        pass a boolean here
        -dull                     just a dull option
    -h, -help                     print help information
    -i, -int <NUMBER>             pass an integer here
    -I, -int64 <64-BIT NUMBER>    pass a 64-bit integer here
    -s, -string <TEXT>            pass a string here (default: "bar")
`,
		},
		{
			desc: "subcommand misuse",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "", root.parg1; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					return nil
				},
				Subcommands: map[string]*cli.Command{
					"foo": {
						Description: "this is foo's description",
						Exec: func(prg cli.Program) error {
							fmt.Fprintln(prg.Stdout(), "testing foo stdout")
							fmt.Fprintln(prg.Stdout(), "testing foo stderr")
							return nil
						},
					},
				},
			},
			args:     []string{"test", "bar"},
			wantCode: 2,
			wantOut:  "",
			wantErr: `test: command provided but not defined: bar
USAGE:
    test [<OPTIONS>] [<COMMAND>]

OPTIONS:
    -h, -help    print help information

COMMANDS:
    foo    this is foo's description
`,
			wantCombined: `test: command provided but not defined: bar
USAGE:
    test [<OPTIONS>] [<COMMAND>]

OPTIONS:
    -h, -help    print help information

COMMANDS:
    foo    this is foo's description
`,
		},
		{
			desc: "subcommand error",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error {
					t := root.T
					if want, got := "", root.parg1; got != want {
						t.Fatalf("want %q, got %q", want, got)
					}
					return nil
				},
				Subcommands: map[string]*cli.Command{
					"foo": {
						Description: "this is foo's description",
						Exec: func(prg cli.Program) error {
							fmt.Fprintln(prg.Stdout(), "testing foo stdout")
							fmt.Fprintln(prg.Stderr(), "testing foo stderr")
							return errors.New("foo has thrown an error")
						},
					},
				},
			},
			args:     []string{"test", "foo"},
			wantCode: 1,
			wantOut:  "testing foo stdout\n",
			wantErr: `test: foo has thrown an error
testing foo stderr
`,
			wantCombined: `test: foo has thrown an error
testing foo stdout
testing foo stderr
`,
		},
		{
			desc: "custom error code",
			entry: &cli.Command{
				Exec: func(_ cli.Program) error { return errors.New("foo") },
			},
			opts: []func(*cli.CLI){
				cli.ErrorCode(0xDEADC0DE),
			},
			args:         []string{"test"},
			wantCode:     0xDEADC0DE,
			wantOut:      "",
			wantErr:      "test: foo\n",
			wantCombined: "test: foo\n",
		},
		{
			desc: "custom misuse code",
			entry: &cli.Command{
				Subcommands: map[string]*cli.Command{
					"foo": new(cli.Command),
				},
			},
			opts: []func(*cli.CLI){
				cli.MisuseCode(0xABADBABE),
			},
			args:     []string{"test", "bar"},
			wantCode: 0xABADBABE,
			wantOut:  "",
			wantErr: `test: command provided but not defined: bar
USAGE:
    test [<OPTIONS>] [<COMMAND>]

OPTIONS:
    -h, -help    print help information

COMMANDS:
    foo
`,
			wantCombined: `test: command provided but not defined: bar
USAGE:
    test [<OPTIONS>] [<COMMAND>]

OPTIONS:
    -h, -help    print help information

COMMANDS:
    foo
`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			defer func() {
				root = testCommand{T: t} // reset root
			}()
			var stdout, stderr, combined strings.Builder
			opts := []func(*cli.CLI){
				cli.Stdout(
					io.MultiWriter(&stdout, &combined),
				),
				cli.Stderr(
					io.MultiWriter(&stderr, &combined),
				),
				cli.HelpDescription("print help information"),
			}
			opts = append(opts, tc.opts...)
			cli := cli.New(
				tc.entry,
				opts...,
			)
			code := cli.ParseAndRun(tc.args)
			if want, got := tc.wantCode, code; got != want {
				t.Fatalf("want %d, got %d", want, got)
			}
			if want, got := tc.wantOut, stdout.String(); got != want {
				t.Fatalf("STDOUT (-want +got):\n%s", cmp.Diff(want, got))
			}
			if want, got := tc.wantErr, stderr.String(); got != want {
				t.Fatalf("STDERR (-want +got):\n%s", cmp.Diff(want, got))
			}
			if want, got := tc.wantCombined, combined.String(); got != want {
				t.Fatalf("STDOUT + STDERR (-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func newDullStr() *string {
	var s string
	return &s
}
