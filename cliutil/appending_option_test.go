package cliutil_test

import (
	"fmt"
	"testing"

	"github.com/gbrlsnchs/cli/cliutil"
	"github.com/google/go-cmp/cmp"
)

func TestAppendingOption(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		testCases := []struct {
			input []string
			want  cliutil.AppendingOption
		}{
			{[]string{""}, cliutil.AppendingOption{}},
			{[]string{"foo"}, cliutil.AppendingOption{"foo"}},
			{[]string{"foo", "bar"}, cliutil.AppendingOption{"foo", "bar"}},
			{[]string{"foo", "bar", "baz"}, cliutil.AppendingOption{"foo", "bar", "baz"}},
			{[]string{"foo", "", "baz"}, cliutil.AppendingOption{"foo", "baz"}},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				ao := make(cliutil.AppendingOption, 0)
				for _, in := range tc.input {
					err := ao.Set(in)
					if want, got := (error)(nil), err; got != want {
						t.Fatalf("want %v, got %v", want, got)
					}
				}
				if want, got := tc.want, ao; !cmp.Equal(got, want) {
					t.Fatalf("(*AppendingOption).Set mismatch (-want +got):\n%s", cmp.Diff(want, got))
				}
			})
		}
	})
	t.Run("String", func(t *testing.T) {
		testCases := []struct {
			ao cliutil.AppendingOption
		}{
			{cliutil.AppendingOption{}},
			{cliutil.AppendingOption{"foo"}},
			{cliutil.AppendingOption{"foo", "bar"}},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				s := tc.ao.String()
				// NOTE: it doesn't matter how this is printed as long it uses the []string formatting
				if want, got := fmt.Sprint([]string(tc.ao)), s; got != want {
					t.Fatalf("want %q, got %q", want, got)
				}
			})
		}
	})
}
