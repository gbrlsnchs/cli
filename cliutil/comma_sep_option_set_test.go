package cliutil_test

import (
	"testing"

	"github.com/gbrlsnchs/cli/cliutil"
	"github.com/google/go-cmp/cmp"
)

func TestCommaSepOptionSet(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		testCases := []struct {
			input string
			want  cliutil.CommaSepOptionSet
		}{
			{"", cliutil.CommaSepOptionSet{}},
			{"foo", cliutil.CommaSepOptionSet{
				"foo": struct{}{},
			}},
			{"foo,bar", cliutil.CommaSepOptionSet{
				"foo": struct{}{},
				"bar": struct{}{},
			}},
			{"foo,bar,baz", cliutil.CommaSepOptionSet{
				"foo": struct{}{},
				"bar": struct{}{},
				"baz": struct{}{},
			}},
			{"foo,,baz", cliutil.CommaSepOptionSet{
				"foo": struct{}{},
				"baz": struct{}{},
			}},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				cs := make(cliutil.CommaSepOptionSet)
				err := cs.Set(tc.input)
				if want, got := (error)(nil), err; got != want {
					t.Fatalf("want %v, got %v", want, got)
				}
				if want, got := tc.want, cs; !cmp.Equal(got, want) {
					t.Fatalf("(*CommaSepOptionSet).Set mismatch (-want +got):\n%s", cmp.Diff(want, got))
				}
			})
		}
		t.Run("reallocate", func(t *testing.T) {
			cs := make(cliutil.CommaSepOptionSet)
			cs.Set("foo,bar")
			cs.Set("baz,qux")
			if want, got := (cliutil.CommaSepOptionSet{"baz": {}, "qux": {}}), cs; !cmp.Equal(got, want) {
				t.Fatalf("(*CommaSepOptionSet).Set doesn't reallocate the underlying set:\n%s", cmp.Diff(want, got))
			}
		})
	})
	t.Run("String", func(t *testing.T) {
		testCases := []struct {
			cs   cliutil.CommaSepOptionSet
			want string
		}{
			{cliutil.CommaSepOptionSet{}, ""},
			{cliutil.CommaSepOptionSet{
				"foo": struct{}{},
			}, "foo"},
			{cliutil.CommaSepOptionSet{
				"foo": struct{}{},
				"bar": struct{}{},
			}, "bar,foo"},
			{cliutil.CommaSepOptionSet{
				"foo": struct{}{},
				"bar": struct{}{},
				"baz": struct{}{},
			}, "bar,baz,foo"},
			{cliutil.CommaSepOptionSet{
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
