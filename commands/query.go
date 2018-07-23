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

	"github.com/chasinglogic/tsk/ql/lexer"
	"github.com/chasinglogic/tsk/ql/parser"
	"github.com/chasinglogic/tsk/task"
	"github.com/spf13/cobra"
)

var query = &cobra.Command{
	Use:     "query",
	Aliases: []string{"q", "s", "search"},
	Short:   "Search and list tasks",
	Run: func(cmd *cobra.Command, args []string) {
		backend, err := config.backend()
		if err != nil {
			fmt.Println("ERROR Unable to load backend:", err)
			os.Exit(1)
		}

		input := strings.Join(args, " ")
		if len(input) == 0 {
			list := backend.Slice()
			printList(list)
			os.Exit(1)
		}

		p := parser.New(lexer.New(input))
		ast := p.Parse()

		if err := p.Error(); err != nil {
			fmt.Println("ERROR parsing query:", err)
			os.Exit(1)
		}

		list, err := backend.Search(ast)
		if err != nil {
			fmt.Println("ERROR searching backend:", err)
			os.Exit(1)
		}

		printList(list)
	},
}

func printList(list []task.Task) {
	for i := range list {
		fmt.Printf("%s \"%s\"\n", list[i].ID, list[i].Title)
	}
}
