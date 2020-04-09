package cliutil

import "strings"

// CommaList is a custom type for comma-separated lists.
type CommaList []string

// Set sets a comma-separated string to a list.
func (cl *CommaList) Set(value string) error {
	if value == "" {
		return nil
	}
	values := strings.Split(value, ",")
	if *cl == nil {
		*cl = make(CommaList, 0, len(values))
	}
	for _, v := range values {
		if v == "" {
			continue
		}
		*cl = append(*cl, v)
	}
	return nil
}

func (cl CommaList) String() string { return strings.Join(cl, ",") }
