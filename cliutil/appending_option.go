package cliutil

import "fmt"

// AppendingOption is a custom option type for repeating flags.
type AppendingOption []string

// Set appends value to the underlying appending option.
func (ao *AppendingOption) Set(value string) error {
	if value == "" {
		return nil
	}
	arr := *ao
	if arr == nil {
		arr = AppendingOption{value}
	}
	arr = append(arr, value)
	*ao = arr
	return nil
}

func (ao AppendingOption) String() string { return fmt.Sprint([]string(ao)) }
