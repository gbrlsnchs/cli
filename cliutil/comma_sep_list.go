package cliutil

import "strings"

// CommaSepList is a custom option type for comma-separated lists.
type CommaSepList []string

// Set sets a comma-separated string to a list.
func (cl *CommaSepList) Set(value string) error {
	if value == "" {
		return nil
	}
	values := strings.Split(value, ",")
	if *cl == nil {
		*cl = make(CommaSepList, 0, len(values))
	}
	for _, v := range values {
		if v == "" {
			continue
		}
		*cl = append(*cl, v)
	}
	return nil
}

func (cl CommaSepList) String() string { return strings.Join(cl, ",") }
