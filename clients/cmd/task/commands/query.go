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

	"github.com/chasinglogic/taskforge/ql/lexer"
	"github.com/chasinglogic/taskforge/ql/parser"
	"github.com/chasinglogic/taskforge/task"
	"github.com/spf13/cobra"
)

func init() {
	query.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}} QUERY...{{end}}{{if .HasAvailableSubCommands}}
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

var query = &cobra.Command{
	Use:     "query",
	Aliases: []string{"q", "s", "search"},
	Short:   "Search and list tasks",
	Run: func(cmd *cobra.Command, args []string) {
		l, err := config.list()
		if err != nil {
			fmt.Println("ERROR Unable to load list:", err)
			os.Exit(1)
		}

		input := strings.Join(args, " ")
		if len(input) == 0 {
			l, err := l.Slice()
			if err != nil {
				fmt.Println("ERROR unable to retrieve tasks:", err)
				os.Exit(1)
			}

			printList(l)
			os.Exit(0)
		}

		p := parser.New(lexer.New(input))
		ast := p.Parse()

		if err := p.Error(); err != nil {
			fmt.Println("ERROR parsing query:", err)
			os.Exit(1)
		}

		result, err := l.Search(ast)
		if err != nil {
			fmt.Println("ERROR searching list:", err)
			os.Exit(1)
		}

		printList(result)
	},
}

func printList(l []task.Task) {
	for i := range l {
		fmt.Printf("%s \"%s\"\n", l[i].ID, l[i].Title)
	}
}
