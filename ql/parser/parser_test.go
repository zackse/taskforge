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
	"encoding/json"
	"reflect"
	"testing"

	"github.com/chasinglogic/taskforge/ql/ast"
	"github.com/chasinglogic/taskforge/ql/lexer"
	"github.com/chasinglogic/taskforge/ql/token"
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
					Operator: token.Token{
						Type:    token.AND,
						Literal: "and",
					},
					Left: ast.StringLiteral{
						Token: token.Token{
							Type:    token.STRING,
							Literal: "milk",
						},
						Value: "milk",
					},
					Right: ast.StringLiteral{
						Token: token.Token{
							Type:    token.STRING,
							Literal: "cookies",
						},
						Value: "cookies",
					},
				},
			},
		},
		{
			name:  "simple all string parse",
			input: "milk -and cookies",
			output: ast.AST{
				Expression: ast.StringLiteral{
					Token: token.Token{
						Type:    token.STRING,
						Literal: "milk and cookies",
					},
					Value: "milk and cookies",
				},
			},
		},
		{
			name:  "complex parse",
			input: "(priority > 5 and title ^ \"take out the trash\") or (context = \"work\" and (priority >= 2 or (\"my little pony\")))",
			output: ast.AST{
				Expression: ast.InfixExpression{
					Operator: token.Token{
						Type:    token.OR,
						Literal: "or",
					},
					Right: ast.InfixExpression{
						Operator: token.Token{
							Type:    token.AND,
							Literal: "and",
						},
						Left: ast.InfixExpression{
							Operator: token.Token{
								Type:    token.EQ,
								Literal: "=",
							},
							Left: ast.StringLiteral{
								Token: token.Token{
									Type:    token.STRING,
									Literal: "context",
								},
								Value: "context",
							},
							Right: ast.StringLiteral{
								Token: token.Token{
									Type:    token.STRING,
									Literal: "work",
								},
								Value: "work",
							},
						},
						Right: ast.InfixExpression{
							Operator: token.Token{
								Type:    token.OR,
								Literal: "or",
							},
							Left: ast.InfixExpression{
								Operator: token.Token{
									Type:    token.GTE,
									Literal: ">=",
								},
								Left: ast.StringLiteral{
									Token: token.Token{
										Type:    token.STRING,
										Literal: "priority",
									},
									Value: "priority",
								},
								Right: ast.NumberLiteral{
									Token: token.Token{
										Type:    token.NUMBER,
										Literal: "2",
									},
									Value: 2.0,
								},
							},
							Right: ast.StringLiteral{
								Token: token.Token{
									Type:    token.STRING,
									Literal: "my little pony",
								},
								Value: "my little pony",
							},
						},
					},
					Left: ast.InfixExpression{
						Operator: token.Token{
							Type:    token.AND,
							Literal: "and",
						},
						Right: ast.InfixExpression{
							Operator: token.Token{
								Type:    token.LIKE,
								Literal: "^",
							},
							Right: ast.StringLiteral{
								Token: token.Token{
									Type:    token.STRING,
									Literal: "take out the trash",
								},
								Value: "take out the trash",
							},
							Left: ast.StringLiteral{
								Token: token.Token{
									Type:    token.STRING,
									Literal: "title",
								},
								Value: "title",
							},
						},
						Left: ast.InfixExpression{
							Operator: token.Token{
								Type:    token.GT,
								Literal: ">",
							},
							Left: ast.StringLiteral{
								Token: token.Token{
									Type:    token.STRING,
									Literal: "priority",
								},
								Value: "priority",
							},
							Right: ast.NumberLiteral{
								Token: token.Token{
									Type:    token.NUMBER,
									Literal: "5",
								},
								Value: 5.0,
							},
						},
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

			if !reflect.DeepEqual(out, test.output) {
				jsn1, _ := json.MarshalIndent(test.output, "", "\t")
				jsn2, _ := json.MarshalIndent(out, "", "\t")
				t.Errorf("Expected %s Got %s", jsn1, jsn2)
			}
		})
	}
}
