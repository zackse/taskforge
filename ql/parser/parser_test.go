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

package parser

import (
	"testing"

	"github.com/chasinglogic/taskforge/ql/ast"
	"github.com/chasinglogic/taskforge/ql/lexer"
	"github.com/chasinglogic/taskforge/ql/token"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		output      ast.AST
		shouldError bool
	}{
		{
			name:  "simple parse",
			input: "milk and cookies",
			output: ast.AST{
				Expression: ast.InfixExpression{
					Operator: token.New("and"),
					Left: ast.StringLiteral{
						Token: token.New("milk"),
						Value: "milk",
					},
					Right: ast.StringLiteral{
						Token: token.New("cookies"),
						Value: "cookies",
					},
				},
			},
		},
		{
			name:  "boolean parse",
			input: "completed = false",
			output: ast.AST{
				Expression: ast.InfixExpression{
					Operator: token.New("="),
					Left: ast.StringLiteral{
						Token: token.New("completed"),
						Value: "completed",
					},
					Right: ast.BooleanLiteral{
						Token: token.New("false"),
						Value: false,
					},
				},
			},
		},
		{
			name:  "simple all string parse",
			input: "milk -and cookies",
			output: ast.AST{
				Expression: ast.StringLiteral{
					Token: token.New("milk and cookies"),
					Value: "milk and cookies",
				},
			},
		},
		{
			name:  "complex parse",
			input: "(priority > 5 and title ^ \"take out the trash\") or (context = \"work\" and (priority >= 2 or (\"my little pony\")))",
			output: ast.AST{
				Expression: ast.InfixExpression{
					Operator: token.New("or"),
					Right: ast.InfixExpression{
						Operator: token.New("and"),
						Left: ast.InfixExpression{
							Operator: token.New("="),
							Left: ast.StringLiteral{
								Token: token.New("context"),
								Value: "context",
							},
							Right: ast.StringLiteral{
								Token: token.New("work"),
								Value: "work",
							},
						},
						Right: ast.InfixExpression{
							Operator: token.New("or"),
							Left: ast.InfixExpression{
								Operator: token.New(">="),
								Left: ast.StringLiteral{
									Token: token.New("priority"),
									Value: "priority",
								},
								Right: ast.NumberLiteral{
									Token: token.New("2"),
									Value: 2.0,
								},
							},
							Right: ast.StringLiteral{
								Token: token.New("my little pony"),
								Value: "my little pony",
							},
						},
					},
					Left: ast.InfixExpression{
						Operator: token.New("and"),
						Right: ast.InfixExpression{
							Operator: token.New("^"),
							Right: ast.StringLiteral{
								Token: token.New("take out the trash"),
								Value: "take out the trash",
							},
							Left: ast.StringLiteral{
								Token: token.New("title"),
								Value: "title",
							},
						},
						Left: ast.InfixExpression{
							Operator: token.New(">"),
							Left: ast.StringLiteral{
								Token: token.New("priority"),
								Value: "priority",
							},
							Right: ast.NumberLiteral{
								Token: token.New("5"),
								Value: 5.0,
							},
						},
					},
				},
			},
		},
		{
			name:  "completed = false",
			input: "completed = false",
			output: ast.AST{
				Expression: ast.InfixExpression{
					Operator: token.New("="),
					Left: ast.StringLiteral{
						Token: token.New("completed"),
						Value: "completed",
					},
					Right: ast.BooleanLiteral{
						Token: token.New("false"),
						Value: false,
					},
				},
			},
		},
		{
			name:  "completed = true",
			input: "completed = true",
			output: ast.AST{
				Expression: ast.InfixExpression{
					Operator: token.New("="),
					Left: ast.StringLiteral{
						Token: token.New("completed"),
						Value: "completed",
					},
					Right: ast.BooleanLiteral{
						Token: token.New("true"),
						Value: true,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lex := lexer.New(test.input)
			par := New(lex)
			out := par.Parse()

			if par.Error() != nil && !test.shouldError {
				t.Errorf("parser errors: %s", par.Error())
				return
			} else if par.Error() == nil && test.shouldError {
				t.Errorf("got no error when should have, output: %s", out)
				return
			} else if par.Error() != nil && test.shouldError {
				return
			}

			require.Equal(t, out, test.output)
		})
	}
}
