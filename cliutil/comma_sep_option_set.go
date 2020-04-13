package cliutil

import (
	"sort"
	"strings"
)

// CommaSepOptionSet is a set for comma-separated options.
type CommaSepOptionSet map[string]struct{}

// Set sets a comma-separated string to a set.
func (cs *CommaSepOptionSet) Set(value string) error {
	if value == "" {
		return nil
	}
	values := strings.Split(value, ",")
	if *cs == nil {
		*cs = make(map[string]struct{}, len(values))
	}
	m := *cs
	for _, v := range values {
		if v == "" {
			continue
		}
		m[v] = struct{}{}
	}
	return nil
}

func (cs CommaSepOptionSet) String() string {
	arr := make([]string, 0, len(cs))
	for k := range cs {
		if k == "" {
			continue
		}
		arr = append(arr, k)
	}
	sort.Strings(arr)
	return strings.Join(arr, ",")
}
