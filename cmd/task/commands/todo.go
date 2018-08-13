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

	"github.com/chasinglogic/taskforge/ql/ast"
	"github.com/chasinglogic/taskforge/ql/token"
	"github.com/spf13/cobra"
)

func init() {
	todo.Flags().StringVarP(&outputFormat, "output", "o", "table", "How to output matching tasks. Options are: table, text, json, csv")
	todo.Flags().BoolVarP(&autoWrapText, "wrap", "w", false,
		`Whether to wrap text or not. For smaller terminals this will improve
	display but for larger terminals this will allow titles to be longer before
	wrapping weirdly`)
}

var todo = &cobra.Command{
	Use:     "todo",
	Aliases: []string{"t"},
	Short:   "Show incomplete tasks",
	Run: func(cmd *cobra.Command, args []string) {
		l, err := config.list()
		if err != nil {
			fmt.Println("ERROR Unable to load list:", err)
			os.Exit(1)
		}

		ast := ast.AST{
			Expression: ast.InfixExpression{
				Operator: token.Token{
					Type: token.EQ,
				},
				Left: ast.StringLiteral{
					Token: token.Token{
						Type:    token.STRING,
						Literal: "completed",
					},
					Value: "completed",
				},
				Right: ast.BooleanLiteral{
					Token: token.Token{
						Type:    token.BOOLEAN,
						Literal: "false",
					},
					Value: false,
				},
			},
		}

		result, err := l.Search(ast)
		if err != nil {
			fmt.Println("ERROR searching list:", err)
			os.Exit(1)
		}

		printList(result)
	},
}
