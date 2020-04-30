package cliutil

import "fmt"

// MultiValueOptionList is a list for options with multiple values.
type MultiValueOptionList []string

// Set appends value to the underlying list.
func (mvl *MultiValueOptionList) Set(value string) error {
	if value == "" {
		return nil
	}
	arr := *mvl
	if arr == nil {
		arr = MultiValueOptionList{value}
	}
	arr = append(arr, value)
	*mvl = arr
	return nil
}

func (mvl MultiValueOptionList) String() string { return fmt.Sprint([]string(mvl)) }
