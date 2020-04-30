package cliutil_test

import (
	"fmt"
	"testing"

	"github.com/gbrlsnchs/cli/cliutil"
	"github.com/google/go-cmp/cmp"
)

func TestMultiValueOptionList(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		testCases := []struct {
			input []string
			want  cliutil.MultiValueOptionList
		}{
			{[]string{""}, cliutil.MultiValueOptionList{}},
			{[]string{"foo"}, cliutil.MultiValueOptionList{"foo"}},
			{[]string{"foo", "bar"}, cliutil.MultiValueOptionList{"foo", "bar"}},
			{[]string{"foo", "bar", "baz"}, cliutil.MultiValueOptionList{"foo", "bar", "baz"}},
			{[]string{"foo", "", "baz"}, cliutil.MultiValueOptionList{"foo", "baz"}},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				mvl := make(cliutil.MultiValueOptionList, 0)
				for _, in := range tc.input {
					err := mvl.Set(in)
					if want, got := (error)(nil), err; got != want {
						t.Fatalf("want %v, got %v", want, got)
					}
				}
				if want, got := tc.want, mvl; !cmp.Equal(got, want) {
					t.Fatalf("(*MultiValueOption).Set mismatch (-want +got):\n%s", cmp.Diff(want, got))
				}
			})
		}
	})
	t.Run("String", func(t *testing.T) {
		testCases := []struct {
			ao cliutil.MultiValueOptionList
		}{
			{cliutil.MultiValueOptionList{}},
			{cliutil.MultiValueOptionList{"foo"}},
			{cliutil.MultiValueOptionList{"foo", "bar"}},
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
