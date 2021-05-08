package command

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
)

// BaseRun contains execution logic for root commands.
func BaseRun(cmd *Command, args []string) {
	if len(args) == 0 {
		cmd.Usage(os.Stderr)
		os.Exit(2)
	}

	if args[0] == "help" {
		cmd.Help(os.Stdout, os.Stderr, args[1:])
		return
	}

	for _, c := range cmd.Commands {
		if args[0] != c.Name {
			continue
		}
		c.Flag.Parse(args[1:])
		if c.Runnable() {
			c.Run(c, c.Flag.Args())
		}
		return
	}

	cmd.Help(os.Stderr, os.Stderr, args[1:])
	os.Exit(2)
}

type Command struct {
	// Name is the command name.
	Name string

	// Arguments is the argument list shown in the 'help' output.
	Arguments string

	// Short is the short description shown in the 'help' output.
	Short string

	// Long is the long message shown in the 'help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag *flag.FlagSet

	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)

	// Commands lists the available commands.
	// The order here is the order in which they are printed by 'help'.
	Commands []*Subcommand
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

var usageTemplate = `{{.Long | trim}}

Usage:

	{{.Name}}{{if (hasRunnable .Commands)}} <command>{{end}}{{if .Arguments}} {{.Arguments}}{{end}}{{if (hasRunnable .Commands)}}

The commands are:
{{range .Commands}}{{if .Runnable}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}{{end}}{{if hasTopic .Commands}}

Additional help topics:
{{range .Commands}}{{if not .Runnable}}
	{{.Name | printf "%-15s"}} {{.Short}}{{end}}{{end}}{{end}}{{if .Commands}}

Use "{{.Name}} help <command>" for more information about a command.{{end}}
`

func (c *Command) Usage(w io.Writer) {
	bw := bufio.NewWriter(w)
	tmpl(bw, usageTemplate, c)
	bw.Flush()
}

var helpTemplate = `{{if .Runnable}}usage: {{.FullName}}{{if .Arguments}} {{.Arguments}}{{end}}

{{end}}{{.Long | trim}}
`

// Help implements the 'help' command.
func (c *Command) Help(out io.Writer, err io.Writer, args []string) {
	if len(c.Commands) == 0 {
		tmpl(out, helpTemplate, struct {
			*Command
			FullName string
		}{
			Command:  c,
			FullName: c.Name,
		})
		return
	}

	if len(args) == 0 {
		c.Usage(out)
		return
	}
	arg := args[0]

	var cmd *Subcommand
	for _, sub := range c.Commands {
		if sub.Name == arg {
			cmd = sub
			break
		}
	}

	if cmd == nil || len(args) > 1 {
		// helpSuccess is the help command using as many args as possible that would succeed.
		helpSuccess := fmt.Sprintf("%s help %s", c.Name, arg)
		fmt.Fprintf(err, "%s help %s: unknown help topic. Run '%s'.\n", c.Name, strings.Join(args, " "), helpSuccess)
		os.Exit(2)
	}

	tmpl(out, helpTemplate, struct {
		*Subcommand
		FullName string
	}{
		Subcommand: cmd,
		FullName:   strings.Join(append([]string{c.Name}, args[0]), " "),
	})
}

type Subcommand struct {
	// Name is the command name.
	Name string

	// Arguments is the argument list shown in the 'help' output.
	Arguments string

	// Short is the short description shown in the 'help' output.
	Short string

	// Long is the long message shown in the 'help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag *flag.FlagSet

	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Subcommand, args []string)
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command.
func (c *Subcommand) Runnable() bool {
	return c.Run != nil
}

// An errWriter wraps a writer, recording whether a write error occurred.
type errWriter struct {
	w   io.Writer
	err error
}

func (w *errWriter) Write(b []byte) (int, error) {
	n, err := w.w.Write(b)
	if err != nil {
		w.err = err
	}
	return n, err
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{
		"trim":        strings.TrimSpace,
		"hasRunnable": hasRunnable,
		"hasTopic":    hasTopic,
	})
	template.Must(t.Parse(text))
	ew := &errWriter{w: w}
	err := t.Execute(ew, data)
	if ew.err != nil {
		// I/O error writing. Ignore write on closed pipe.
		if strings.Contains(ew.err.Error(), "pipe") {
			os.Exit(1)
		}
		log.Printf("writing output: %v", ew.err)
		os.Exit(1)
	}
	if err != nil {
		panic(err)
	}
}

func hasRunnable(cmds []*Subcommand) bool {
	for _, cmd := range cmds {
		if cmd.Runnable() {
			return true
		}
	}
	return false
}

func hasTopic(cmds []*Subcommand) bool {
	for _, cmd := range cmds {
		if !cmd.Runnable() {
			return true
		}
	}
	return false
}
