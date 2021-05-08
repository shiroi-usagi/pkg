package command

import (
	"bytes"
	"testing"
)

func TestCommand_Usage(t *testing.T) {
	tests := []struct {
		name   string
		target *Command
		wantW  string
	}{
		{
			name: "empty",
			target: &Command{
				Name: "test",
				Long: `Long description of test command.`,
			},
			wantW: `Long description of test command.

Usage:

	test
`,
		},
		{
			name: "arguments",
			target: &Command{
				Name:      "test",
				Arguments: "[arguments]",
				Long:      `Long description of test command.`,
			},
			wantW: `Long description of test command.

Usage:

	test [arguments]
`,
		},
		{
			name: "subcommand",
			target: &Command{
				Name:      "test",
				Arguments: "[arguments]",
				Long:      `Long description of test command.`,
				Commands: []*Subcommand{
					{
						Name:  "test-subcommand",
						Short: "short",
						Long:  "Long description",
						Run:   func(cmd *Subcommand, args []string) {},
					},
				},
			},
			wantW: `Long description of test command.

Usage:

	test <command> [arguments]

The commands are:

	test-subcommand short

Use "test help <command>" for more information about a command.
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.target
			w := &bytes.Buffer{}
			c.Usage(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Usage() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestCommand_Help(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		target  *Command
		args    args
		wantOut string
		wantErr string
	}{
		{
			name: "empty",
			target: &Command{
				Name: "test",
				Long: `Long description of test command.`,
			},
			args: args{[]string{}},
			wantOut: `Long description of test command.
`,
		},
		{
			name: "runnable",
			target: &Command{
				Name: "test",
				Long: `Long description of test command.`,
				Run:  func(_ *Command, _ []string) {},
			},
			args: args{[]string{}},
			wantOut: `usage: test

Long description of test command.
`,
		},
		{
			name: "with subcommand",
			target: &Command{
				Name: "test",
				Long: `Long description of test command.`,
				Run:  func(cmd *Command, args []string) {},
				Commands: []*Subcommand{
					{
						Name:  "subcommand",
						Short: "shor desc",
						Run:   func(_ *Subcommand, _ []string) {},
					},
				},
			},
			args: args{[]string{}},
			wantOut: `Long description of test command.

Usage:

	test <command>

The commands are:

	subcommand  shor desc

Use "test help <command>" for more information about a command.
`,
		},
		{
			name: "with topic",
			target: &Command{
				Name: "test",
				Long: `Long description of test command.`,
				Run:  func(cmd *Command, args []string) {},
				Commands: []*Subcommand{
					{
						Name:  "subcommand",
						Short: "shor desc",
					},
				},
			},
			args: args{[]string{}},
			wantOut: `Long description of test command.

Usage:

	test

Additional help topics:

	subcommand      shor desc

Use "test help <command>" for more information about a command.
`,
		},
		{
			name: "subcommand",
			target: &Command{
				Name: "test",
				Long: `Long description of test command.`,
				Commands: []*Subcommand{
					{
						Name: "subcommand",
						Long: `Long description of subcommand.`,
						Run:  func(_ *Subcommand, _ []string) {},
					},
				},
			},
			args: args{[]string{"subcommand"}},
			wantOut: `usage: test subcommand

Long description of subcommand.
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.target
			out := &bytes.Buffer{}
			err := &bytes.Buffer{}
			c.Help(out, err, tt.args.args)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("Help() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
			if gotErr := err.String(); gotErr != tt.wantErr {
				t.Errorf("Help() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
