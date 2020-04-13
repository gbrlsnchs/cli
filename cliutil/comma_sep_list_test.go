package cliutil_test

import (
	"testing"

	"github.com/gbrlsnchs/cli/cliutil"
	"github.com/google/go-cmp/cmp"
)

func TestCommaSepList(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		testCases := []struct {
			input string
			want  cliutil.CommaSepList
		}{
			{"", cliutil.CommaSepList{}},
			{"foo", cliutil.CommaSepList{"foo"}},
			{"foo,bar", cliutil.CommaSepList{"foo", "bar"}},
			{"foo,bar,baz", cliutil.CommaSepList{"foo", "bar", "baz"}},
			{"foo,,baz", cliutil.CommaSepList{"foo", "baz"}},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				cs := make(cliutil.CommaSepList, 0)
				err := cs.Set(tc.input)
				if want, got := (error)(nil), err; got != want {
					t.Fatalf("want %v, got %v", want, got)
				}
				if want, got := tc.want, cs; !cmp.Equal(got, want) {
					t.Fatalf("(*CommaSepList).Set mismatch (-want +got):\n%s", cmp.Diff(want, got))
				}
			})
		}
	})
	t.Run("String", func(t *testing.T) {
		testCases := []struct {
			cs   cliutil.CommaSepList
			want string
		}{
			{cliutil.CommaSepList{}, ""},
			{cliutil.CommaSepList{"foo"}, "foo"},
			{cliutil.CommaSepList{"foo", "bar"}, "foo,bar"},
			{cliutil.CommaSepList{"foo", "bar", "baz"}, "foo,bar,baz"},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				s := tc.cs.String()
				if want, got := tc.want, s; got != want {
					t.Fatalf("want %q, got %q", want, got)
				}
			})
		}
	})
}
