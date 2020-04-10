package clitest_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/gbrlsnchs/cli/clitest"
	"github.com/google/go-cmp/cmp"
)

func TestProgram(t *testing.T) {
	t.Run("Name", func(t *testing.T) {
		prg := clitest.NewProgram("test")
		if want, got := "test", prg.Name(); got != want {
			t.Fatalf("want %q, got %q", want, got)
		}
	})

	type output int
	const (
		stdout output = iota
		stderr
	)
	type input struct {
		o   output
		txt string
	}

	testCases := []struct {
		inputs        []input
		wantOutput    string
		wantErrOutput string
		wantCombined  string
	}{
		{
			inputs: []input{
				{stdout, "foo"},
				{stderr, "bar"},
			},
			wantOutput:    "foo\n",
			wantErrOutput: "bar\n",
			wantCombined:  "foo\nbar\n",
		},
		{
			inputs: []input{
				{stderr, "foo"},
				{stdout, "bar"},
			},
			wantOutput:    "bar\n",
			wantErrOutput: "foo\n",
			wantCombined:  "foo\nbar\n",
		},
	}
	for _, tc := range testCases {
		prg := clitest.NewProgram("test")
		for _, in := range tc.inputs {
			var w io.Writer
			switch in.o {
			case stdout:
				w = prg.Stdout()
			case stderr:
				w = prg.Stderr()
			}
			fmt.Fprintln(w, in.txt)
		}
		t.Run("Output", func(t *testing.T) {
			t.Run("", func(t *testing.T) {
				if want, got := tc.wantOutput, prg.Output(); got != want {
					t.Fatalf("Program.Output mismatch (-want +got):\n%s", cmp.Diff(want, got))
				}
			})
		})
		t.Run("ErrOutput", func(t *testing.T) {
			t.Run("", func(t *testing.T) {
				if want, got := tc.wantErrOutput, prg.ErrOutput(); got != want {
					t.Fatalf("Program.ErrOutput mismatch (-want +got):\n%s", cmp.Diff(want, got))
				}
			})
		})
		t.Run("CombinedOutput", func(t *testing.T) {
			t.Run("", func(t *testing.T) {
				if want, got := tc.wantCombined, prg.CombinedOutput(); got != want {
					t.Fatalf("Program.CombinedOutput mismatch (-want +got):\n%s", cmp.Diff(want, got))
				}
			})
		})
	}
}
