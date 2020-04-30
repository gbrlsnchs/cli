package cliutil

import "fmt"

// MultiValueOptionSet is a set for options with multiple values.
type MultiValueOptionSet map[string]struct{}

// Set adds value to the underlying set.
func (mvs *MultiValueOptionSet) Set(value string) error {
	if value == "" {
		return nil
	}
	if *mvs == nil {
		*mvs = make(map[string]struct{})
	}
	(*mvs)[value] = struct{}{}
	return nil
}

func (mvs MultiValueOptionSet) String() string { return fmt.Sprint(map[string]struct{}(mvs)) }
