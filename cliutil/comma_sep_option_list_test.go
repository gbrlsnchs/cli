package cliutil_test

import (
	"testing"

	"github.com/gbrlsnchs/cli/cliutil"
	"github.com/google/go-cmp/cmp"
)

func TestCommaSepOptionList(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		testCases := []struct {
			input string
			want  cliutil.CommaSepOptionList
		}{
			{"", cliutil.CommaSepOptionList{}},
			{"foo", cliutil.CommaSepOptionList{"foo"}},
			{"foo,bar", cliutil.CommaSepOptionList{"foo", "bar"}},
			{"foo,bar,baz", cliutil.CommaSepOptionList{"foo", "bar", "baz"}},
			{"foo,,baz", cliutil.CommaSepOptionList{"foo", "baz"}},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				cs := make(cliutil.CommaSepOptionList, 0)
				err := cs.Set(tc.input)
				if want, got := (error)(nil), err; got != want {
					t.Fatalf("want %v, got %v", want, got)
				}
				if want, got := tc.want, cs; !cmp.Equal(got, want) {
					t.Fatalf("(*CommaSepList).Set mismatch (-want +got):\n%s", cmp.Diff(want, got))
				}
			})
		}
		t.Run("reallocate", func(t *testing.T) {
			cl := make(cliutil.CommaSepOptionList, 0)
			cl.Set("foo,bar")
			cl.Set("baz,qux")
			if want, got := (cliutil.CommaSepOptionList{"baz", "qux"}), cl; !cmp.Equal(got, want) {
				t.Fatalf("(*CommaSepOptionSet).Set doesn't reallocate the underlying set:\n%s", cmp.Diff(want, got))
			}
		})
	})
	t.Run("String", func(t *testing.T) {
		testCases := []struct {
			cs   cliutil.CommaSepOptionList
			want string
		}{
			{cliutil.CommaSepOptionList{}, ""},
			{cliutil.CommaSepOptionList{"foo"}, "foo"},
			{cliutil.CommaSepOptionList{"foo", "bar"}, "foo,bar"},
			{cliutil.CommaSepOptionList{"foo", "bar", "baz"}, "foo,bar,baz"},
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
