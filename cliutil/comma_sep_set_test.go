package cliutil_test

import (
	"testing"

	"github.com/gbrlsnchs/cli/cliutil"
	"github.com/google/go-cmp/cmp"
)

func TestCommaSepSet(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		testCases := []struct {
			input string
			want  cliutil.CommaSepSet
		}{
			{"", cliutil.CommaSepSet{}},
			{"foo", cliutil.CommaSepSet{
				"foo": struct{}{},
			}},
			{"foo,bar", cliutil.CommaSepSet{
				"foo": struct{}{},
				"bar": struct{}{},
			}},
			{"foo,bar,baz", cliutil.CommaSepSet{
				"foo": struct{}{},
				"bar": struct{}{},
				"baz": struct{}{},
			}},
			{"foo,,baz", cliutil.CommaSepSet{
				"foo": struct{}{},
				"baz": struct{}{},
			}},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				cs := make(cliutil.CommaSepSet)
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
			cs   cliutil.CommaSepSet
			want string
		}{
			{cliutil.CommaSepSet{}, ""},
			{cliutil.CommaSepSet{
				"foo": struct{}{},
			}, "foo"},
			{cliutil.CommaSepSet{
				"foo": struct{}{},
				"bar": struct{}{},
			}, "bar,foo"},
			{cliutil.CommaSepSet{
				"foo": struct{}{},
				"bar": struct{}{},
				"baz": struct{}{},
			}, "bar,baz,foo"},
			{cliutil.CommaSepSet{
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
