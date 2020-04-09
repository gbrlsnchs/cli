package cli

import (
	"flag"
	"fmt"
	"io"
)

// Option is a type that is able to define its flags to a flag set
// and also print its own documentation.
type Option interface {
	Define(f *flag.FlagSet, name string)
	WriteDoc(w io.Writer, name string)
}

// OptionDetails are common fields for an option, which are its details.
type OptionDetails struct {
	Description string
	Short       byte
	ArgLabel    string
}

// WriteDoc writes to w a flag's description in a pretty way.
func (ff OptionDetails) WriteDoc(w io.Writer, name string) {
	if ff.Short == 0 {
		fmt.Fprint(w, "    ")
	} else {
		fmt.Fprintf(w, "-%s, ", string(ff.Short))
	}
	fmt.Fprintf(w, "-%s", name)
	if ff.ArgLabel != "" {
		fmt.Fprintf(w, " <%s>", ff.ArgLabel)
	}
	fmt.Fprintf(w, "\t%s", ff.Description)
}

func (ff OptionDetails) defineShort(f *flag.FlagSet, name string) {
	if ff.Short == 0 {
		return
	}
	fg := f.Lookup(name)
	if fg == nil {
		return
	}
	f.Var(fg.Value, string(ff.Short), ff.Description)
}

// BoolOption represents a boolean flag.
type BoolOption struct {
	OptionDetails
	DefValue  bool
	Recipient *bool
}

// Define implements Option by defining a boolean flag to f.
func (fg BoolOption) Define(f *flag.FlagSet, name string) {
	f.BoolVar(fg.Recipient, name, fg.DefValue, fg.Description)
	fg.defineShort(f, name)
}

// StringOption represents a string flag.
type StringOption struct {
	OptionDetails
	DefValue  string
	Recipient *string
}

// Define implements Option by defining a string flag to f.
func (fg StringOption) Define(f *flag.FlagSet, name string) {
	f.StringVar(fg.Recipient, name, fg.DefValue, fg.Description)
	fg.defineShort(f, name)
}

// WriteDoc writes the standard flag documentation and also the default
// value when it's not an empty string to w.
func (fg StringOption) WriteDoc(w io.Writer, name string) {
	fg.OptionDetails.WriteDoc(w, name)
	if fg.DefValue == "" {
		return
	}
	fmt.Fprintf(w, " (default: %q)", fg.DefValue)
}

// IntOption represents an integer flag.
type IntOption struct {
	OptionDetails
	DefValue  int
	Recipient *int
}

// Define implements Option by defining an integer flag to f.
func (fg IntOption) Define(f *flag.FlagSet, name string) {
	f.IntVar(fg.Recipient, name, fg.DefValue, fg.Description)
	fg.defineShort(f, name)
}

// Int64Option represents a 64-bit integer flag.
type Int64Option struct {
	OptionDetails
	DefValue  int64
	Recipient *int64
}

// Define implements Option by defining a 64-bit integer flag to f.
func (fg Int64Option) Define(f *flag.FlagSet, name string) {
	f.Int64Var(fg.Recipient, name, fg.DefValue, fg.Description)
	fg.defineShort(f, name)
}

// VarOption represents a flag that implements flag.Value.
type VarOption struct {
	OptionDetails
	Recipient flag.Value
}

// Define implements Option by defining a flag that implements flag.Value to f.
func (fg VarOption) Define(f *flag.FlagSet, name string) {
	f.Var(fg.Recipient, name, fg.Description)
	fg.defineShort(f, name)
}
