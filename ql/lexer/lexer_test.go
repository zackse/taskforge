package lexer

import (
	"reflect"
	"testing"

	"github.com/chasinglogic/tsk/ql/token"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token.Token
	}{
		{
			name:  "simple lex",
			input: "milk and cookies",
			expected: []token.Token{
				{
					Type:    token.STRING,
					Literal: "milk",
				},
				{
					Type:    token.AND,
					Literal: "and",
				},
				{
					Type:    token.STRING,
					Literal: "cookies",
				},
			},
		},
		{
			name:  "keyword excaped lex",
			input: "milk -and cookies",
			expected: []token.Token{
				{
					Type:    token.STRING,
					Literal: "milk",
				},
				{
					Type:    token.STRING,
					Literal: "and",
				},
				{
					Type:    token.STRING,
					Literal: "cookies",
				},
			},
		},
		{
			name:  "complicated lex",
			input: "(priority > 5 and title ^ \"take out the trash\") or (context = \"work\" and (priority >= 2 or (\"my little pony\")))",
			expected: []token.Token{
				{
					Type:    token.LPAREN,
					Literal: "(",
				},
				{
					Type:    token.STRING,
					Literal: "priority",
				},
				{
					Type:    token.GT,
					Literal: ">",
				},
				{
					Type:    token.NUMBER,
					Literal: "5",
				},
				{
					Type:    token.AND,
					Literal: "and",
				},
				{
					Type:    token.STRING,
					Literal: "title",
				},
				{
					Type:    token.LIKE,
					Literal: "^",
				},
				{
					Type:    token.STRING,
					Literal: "take out the trash",
				},
				{
					Type:    token.RPAREN,
					Literal: ")",
				},
				{
					Type:    token.OR,
					Literal: "or",
				},
				{
					Type:    token.LPAREN,
					Literal: "(",
				},
				{
					Type:    token.STRING,
					Literal: "context",
				},
				{
					Type:    token.EQ,
					Literal: "=",
				},
				{
					Type:    token.STRING,
					Literal: "work",
				},
				{
					Type:    token.AND,
					Literal: "and",
				},
				{
					Type:    token.LPAREN,
					Literal: "(",
				},
				{
					Type:    token.STRING,
					Literal: "priority",
				},
				{
					Type:    token.GTE,
					Literal: ">=",
				},
				{
					Type:    token.NUMBER,
					Literal: "2",
				},
				{
					Type:    token.OR,
					Literal: "or",
				},
				{
					Type:    token.LPAREN,
					Literal: "(",
				},
				{
					Type:    token.STRING,
					Literal: "my little pony",
				},
				{
					Type:    token.RPAREN,
					Literal: ")",
				},
				{
					Type:    token.RPAREN,
					Literal: ")",
				},
				{
					Type:    token.RPAREN,
					Literal: ")",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lex := New(test.input)

			for i := range test.expected {
				tok := lex.NextToken()
				if !reflect.DeepEqual(test.expected[i], tok) {
					t.Errorf("Expected: %v Got: %v", test.expected[i], tok)
				}
			}
		})
	}
}
