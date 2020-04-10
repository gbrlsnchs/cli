package cli

import (
	"fmt"
	"sort"
	"strings"
)

const defaultLeftPadding = 4

var usageTemplate = fmt.Sprintf(
	"{{if .Command.Description}}{{.Command.Description}}\n{{end}}" +
		"USAGE:\n" +
		"\t{{.Command.Usage}}\n" +
		"\nOPTIONS:\n" +
		"{{range .Options}}\t{{.Description}}\n{{end}}" +
		"{{if .Command.Subcommands}}\nCOMMANDS:\n{{end}}" +
		"{{range $key, $value := .Command.Subcommands}}\t{{$key}}{{if .Description}}\t{{end}}{{.Description}}\n{{end}}",
)

type commandTemplate struct {
	Command *Command
	Name    string
	Options []optionTemplate
}

type optionTemplate struct {
	Description string
}

func buildTemplate(name string, c *Command) commandTemplate {
	tmpl := commandTemplate{
		Command: c,
		Name:    name,
		Options: make([]optionTemplate, 0, len(c.Options)),
	}
	flags := make([]string, 0, len(c.Options))
	for k := range c.Options {
		flags = append(flags, k)
	}
	sort.Strings(flags)
	bd := new(strings.Builder)
	for _, name := range flags {
		fg := c.Options[name]
		fg.WriteDoc(bd, name)
		tmpl.Options = append(tmpl.Options, optionTemplate{bd.String()})
		bd.Reset()
	}
	return tmpl
}
