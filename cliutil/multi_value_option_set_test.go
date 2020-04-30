package cliutil_test

import (
	"fmt"
	"testing"

	"github.com/gbrlsnchs/cli/cliutil"
	"github.com/google/go-cmp/cmp"
)

func TestMultiValueOptionSet(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		testCases := []struct {
			input []string
			want  cliutil.MultiValueOptionSet
		}{
			{[]string{""}, cliutil.MultiValueOptionSet{}},
			{[]string{"foo"}, cliutil.MultiValueOptionSet{"foo": {}}},
			{[]string{"foo", "bar"}, cliutil.MultiValueOptionSet{"foo": {}, "bar": {}}},
			{[]string{"foo", "bar", "baz"}, cliutil.MultiValueOptionSet{"foo": {}, "bar": {}, "baz": {}}},
			{[]string{"foo", "", "baz"}, cliutil.MultiValueOptionSet{"foo": {}, "baz": {}}},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				mvs := make(cliutil.MultiValueOptionSet)
				for _, in := range tc.input {
					err := mvs.Set(in)
					if want, got := (error)(nil), err; got != want {
						t.Fatalf("want %v, got %v", want, got)
					}
				}
				if want, got := tc.want, mvs; !cmp.Equal(got, want) {
					t.Fatalf("(*MultiValueOption).Set mismatch (-want +got):\n%s", cmp.Diff(want, got))
				}
			})
		}
	})
	t.Run("String", func(t *testing.T) {
		testCases := []struct {
			mvs cliutil.MultiValueOptionSet
		}{
			{cliutil.MultiValueOptionSet{}},
			{cliutil.MultiValueOptionSet{"foo": {}}},
			{cliutil.MultiValueOptionSet{"foo": {}, "bar": {}}},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				s := tc.mvs.String()
				// NOTE: it doesn't matter how this is printed as long it uses the []string formatting
				if want, got := fmt.Sprint(map[string]struct{}(tc.mvs)), s; got != want {
					t.Fatalf("want %q, got %q", want, got)
				}
			})
		}
	})
}
