package cli

import (
	"fmt"
	"io"
)

// Arg is an interface for a positional argument.
// It can append an argument to an argument list and can write its
// own documentation to show on a command's usage instructions.
type Arg interface {
	AppendTo(a *ArgList)
	WriteDoc(w io.Writer)
}

// StringArg is the most common type of argument, a simple string.
type StringArg struct {
	Label     string  // Label is for documentation purposes.
	Required  bool    // Required triggers an error when the argument is not provided.
	Recipient *string // Recipient is the pointer to have the value set to.
	Next      Arg     // Next is the next positional argument.
}

// AppendTo appends the argument and recursively appends chained arguments.
func (arg StringArg) AppendTo(a *ArgList) {
	a.Append(arg.Label, (*strValue)(arg.Recipient), arg.Required, false)
	if next := arg.Next; next != nil {
		next.AppendTo(a)
	}
}

// WriteDoc writes the argument's instruction to w.
func (arg StringArg) WriteDoc(w io.Writer) {
	fmt.Fprintf(w, " ")
	if !arg.Required {
		fmt.Fprint(w, "[")
		defer fmt.Fprint(w, "]")
	}
	fmt.Fprintf(w, "<%s>", arg.Label)
	if arg.Next != nil {
		arg.Next.WriteDoc(w)
	}
}

// RepeatingArg is a repeating argument. It can be empty when not required,
// or must occur one or more times when required.
type RepeatingArg struct {
	Label     string    // Label is for documentation purposes.
	Required  bool      // Required means one or more occurrences must happen.
	Recipient *[]string // Recipient is the pointer that will receive the parsed args.
}

// AppendTo appends the argument as the last one in the list.
func (arg RepeatingArg) AppendTo(a *ArgList) {
	a.Append(arg.Label, (*listValue)(arg.Recipient), arg.Required, true)
}

// WriteDoc writes the argument's instruction to w.
func (arg RepeatingArg) WriteDoc(w io.Writer) {
	fmt.Fprintf(w, " ")
	if !arg.Required {
		fmt.Fprint(w, "[")
		defer fmt.Fprint(w, "]")
	}
	fmt.Fprintf(w, "<%s>...", arg.Label)
}

type argument struct {
	name     string
	required bool
	repeat   bool
	value    ArgValue
}

// ArgList is an argument list that holds all arguments set by a command.
type ArgList struct {
	args []argument
}

// Append appends an argument to itself.
func (a *ArgList) Append(name string, v ArgValue, required, repeat bool) {
	a.args = append(a.args, argument{
		name:     name,
		required: required,
		repeat:   repeat,
		value:    v,
	})
}

func (a *ArgList) missing(args []string) *argument {
	var (
		stk     = a.args
		current int
	)
	for current = 0; current < len(args) && current < len(a.args); current++ {
	}
	stk = stk[current:]
	args = args[current:]
	for i, arg := range stk {
		if !arg.required {
			return nil
		}
		count := i + 1
		if count > len(args) {
			return &stk[i]
		}
		if arg.repeat { // if arg repeats, next repeating args are optional
			return nil
		}
	}
	return nil
}

func (a *ArgList) parse(args []string) error {
	for i := 0; i < len(args) && i < len(a.args); i++ {
		arg := a.args[i]
		if arg.repeat {
			return arg.value.Set(args[i:])
		}
		if err := arg.value.Set([]string{args[i]}); err != nil {
			return err
		}
	}
	return nil
}

// ArgValue is an interface to wrap the parsed arguments.
// Non-repeating arguments get a single argument slice passed to Set.
type ArgValue interface {
	Set([]string) error
}

type strValue string

func (sv *strValue) Set(v []string) error {
	*sv = strValue(v[0])
	return nil
}

type listValue []string

func (lv *listValue) Set(v []string) error {
	*lv = listValue(v)
	return nil
}
