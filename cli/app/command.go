package app

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"strings"
)

// Command type of command on ci.
type Command struct {
	Name        string
	Run         func(cmd *Command, arg []string) int
	Info        template.HTML
	Usage       template.HTML
	Flag        flag.FlagSet
	CustomFlags bool
}

// ShowUsage print usage of the command.
func (c *Command) ShowUsage() {
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(string(c.Usage)))
	os.Exit(2)
}

// IsRunable check is this command is can be run.
func (c *Command) IsRunable() bool {
	return c.Run != nil
}
