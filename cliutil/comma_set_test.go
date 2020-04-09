package cliutil_test

import (
	"testing"

	"github.com/gbrlsnchs/cli/cliutil"
	"github.com/google/go-cmp/cmp"
)

func TestCommaSet(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		testCases := []struct {
			input string
			want  cliutil.CommaSet
		}{
			{"", cliutil.CommaSet{}},
			{"foo", cliutil.CommaSet{
				"foo": struct{}{},
			}},
			{"foo,bar", cliutil.CommaSet{
				"foo": struct{}{},
				"bar": struct{}{},
			}},
			{"foo,bar,baz", cliutil.CommaSet{
				"foo": struct{}{},
				"bar": struct{}{},
				"baz": struct{}{},
			}},
			{"foo,,baz", cliutil.CommaSet{
				"foo": struct{}{},
				"baz": struct{}{},
			}},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				cs := make(cliutil.CommaSet)
				err := cs.Set(tc.input)
				if want, got := (error)(nil), err; got != want {
					t.Fatalf("want %v, got %v", want, got)
				}
				if want, got := tc.want, cs; !cmp.Equal(got, want) {
					t.Fatalf("(*CommaSet).Set mismatch (-want +got):\n%s", cmp.Diff(want, got))
				}
			})
		}
	})
	t.Run("String", func(t *testing.T) {
		testCases := []struct {
			cs   cliutil.CommaSet
			want string
		}{
			{cliutil.CommaSet{}, ""},
			{cliutil.CommaSet{
				"foo": struct{}{},
			}, "foo"},
			{cliutil.CommaSet{
				"foo": struct{}{},
				"bar": struct{}{},
			}, "bar,foo"},
			{cliutil.CommaSet{
				"foo": struct{}{},
				"bar": struct{}{},
				"baz": struct{}{},
			}, "bar,baz,foo"},
			{cliutil.CommaSet{
				"foo": struct{}{},
				"":    struct{}{},
				"baz": struct{}{},
			}, "baz,foo"},
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
