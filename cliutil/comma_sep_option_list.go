package cliutil

import "strings"

// CommaSepOptionList is a list for comma-separated options.
type CommaSepOptionList []string

// Set sets a comma-separated string to a list.
func (cl *CommaSepOptionList) Set(value string) error {
	if value == "" {
		return nil
	}
	values := strings.Split(value, ",")
	if *cl == nil {
		*cl = make(CommaSepOptionList, 0, len(values))
	}
	for _, v := range values {
		if v == "" {
			continue
		}
		*cl = append(*cl, v)
	}
	return nil
}

func (cl CommaSepOptionList) String() string { return strings.Join(cl, ",") }
