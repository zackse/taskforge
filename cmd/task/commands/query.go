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
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/chasinglogic/taskforge/ql/lexer"
	"github.com/chasinglogic/taskforge/ql/parser"
	"github.com/chasinglogic/taskforge/task"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	autoWrapText bool
	outputFormat string
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

	query.Flags().StringVarP(&outputFormat, "output", "o", "table", "How to output matching tasks. Options are: table, text, json, csv")
	query.Flags().BoolVarP(&autoWrapText, "wrap", "w", false,
		`Whether to wrap text or not. For smaller terminals this will improve
	display but for larger terminals this will allow titles to be longer before
	wrapping weirdly`)
}

var query = &cobra.Command{
	Use:     "query",
	Aliases: []string{"q", "s", "search", "list"},
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

func printList(taskList []task.Task) {
	switch outputFormat {
	case "text":
		printText(taskList)
	case "csv":
		printCSV(taskList)
	case "json":
		printJSON(taskList)
	case "table":
		printTable(taskList)
	default:
		fmt.Printf("%s is not a valid output format, defaulting to table.", outputFormat)
		printTable(taskList)
	}
}

func printTable(taskList []task.Task) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title"})
	table.SetRowLine(true)
	table.SetAutoWrapText(autoWrapText)

	for _, t := range taskList {
		table.Append([]string{t.ID, t.Title})
	}

	table.Render()
}

func printJSON(taskList []task.Task) {
	jsn, err := json.MarshalIndent(taskList, "", "\t")
	if err != nil {
		fmt.Println("ERROR unable to marshal json:", err)
	}

	fmt.Println(string(jsn))
}

func printCSV(taskList []task.Task) {
	// print headers
	fmt.Println("ID,Title,Context,Priority,Body,CreatedDate,CompletedDate")

	// print rows
	for _, t := range taskList {
		fmt.Printf(
			//id,title,context,priority,body,createdDate,
			"%s,%s,%s,%.1f,%s,%s,",
			t.ID,
			t.Title,
			t.Context,
			t.Priority,
			t.Body,
			t.CreatedDate.String(),
		)

		// if completed print date
		if t.IsCompleted() {
			fmt.Printf("%s\n", t.CompletedDate.String())
		} else {
			fmt.Printf("\n")
		}
	}
}

func printText(taskList []task.Task) {
	for _, t := range taskList {
		fmt.Println(t.ID, t.Title)
	}
}
