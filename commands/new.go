// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/chasinglogic/tsk/task"
	"github.com/spf13/cobra"
)

var (
	context  string
	body     string
	fromFile string
	priority float64
)

func init() {
	new.Flags().StringVarP(&context, "context", "c", "default",
		"the context which to create this task in")
	new.Flags().StringVarP(&body, "body", "b", "",
		"the text body to give this task")
	new.Flags().StringVarP(&fromFile, "from-file", "", "",
		"a yaml or csv file of tasks to create")
	new.Flags().Float64VarP(&priority, "priority", "p", 0.0,
		"priority with which to give this task")

	new.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}} TITLE...{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}
Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
}

var new = &cobra.Command{
	Use:     "new",
	Aliases: []string{"n", "create"},
	Short:   "Create a new task",
	Run: func(cmd *cobra.Command, args []string) {
		// if fromFile != "" {
		// 	createTasksFromFile()
		// 	return
		// }

		title := strings.Join(args, " ")
		t := task.New(title)
		t.Context = context
		t.Body = body
		t.Priority = priority

		backend, err := config.backend()
		if err != nil {
			fmt.Println("ERROR Unable to load backend:", err)
			os.Exit(1)
		}

		err = backend.Add(t)
		if err != nil {
			fmt.Println("ERROR Unable to add task:", err)
		}

		err = backend.Save()
		if err != nil {
			fmt.Println("ERROR Unable to save to backend:", err)
		}
	},
}
