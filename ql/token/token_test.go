package token

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewToken(t *testing.T) {
	tests := []struct {
		name    string
		literal string
		output  Token
	}{
		{
			name:    "bool true",
			literal: "true",
			output: Token{
				Type:    BOOLEAN,
				Literal: "true",
			},
		},
		{
			name:    "bool True",
			literal: "True",
			output: Token{
				Type:    BOOLEAN,
				Literal: "True",
			},
		},
		{
			name:    "bool false",
			literal: "false",
			output: Token{
				Type:    BOOLEAN,
				Literal: "false",
			},
		},
		{
			name:    "bool False",
			literal: "False",
			output: Token{
				Type:    BOOLEAN,
				Literal: "False",
			},
		},
		{
			name:    "LPAREN (",
			literal: "(",
			output: Token{
				Type:    LPAREN,
				Literal: "(",
			},
		},
		{
			name:    "RPAREN )",
			literal: ")",
			output: Token{
				Type:    RPAREN,
				Literal: ")",
			},
		},
		{
			name:    "GT >",
			literal: ">",
			output: Token{
				Type:    GT,
				Literal: ">",
			},
		},
		{
			name:    "GTE >=",
			literal: ">=",
			output: Token{
				Type:    GTE,
				Literal: ">=",
			},
		},
		{
			name:    "LT <",
			literal: "<",
			output: Token{
				Type:    LT,
				Literal: "<",
			},
		},
		{
			name:    "LTE <=",
			literal: "<=",
			output: Token{
				Type:    LTE,
				Literal: "<=",
			},
		},
		{
			name:    "EQ =",
			literal: "=",
			output: Token{
				Type:    EQ,
				Literal: "=",
			},
		},
		{
			name:    "NE !=",
			literal: "!=",
			output: Token{
				Type:    NE,
				Literal: "!=",
			},
		},
		{
			name:    "NE ^=",
			literal: "^=",
			output: Token{
				Type:    NE,
				Literal: "^=",
			},
		},
		{
			name:    "LIKE ^",
			literal: "^",
			output: Token{
				Type:    LIKE,
				Literal: "^",
			},
		},
		{
			name:    "LIKE ~",
			literal: "~",
			output: Token{
				Type:    LIKE,
				Literal: "~",
			},
		},
		{
			name:    "NLIKE !~",
			literal: "!~",
			output: Token{
				Type:    NLIKE,
				Literal: "!~",
			},
		},
		{
			name:    "NLIKE ^^",
			literal: "^^",
			output: Token{
				Type:    NLIKE,
				Literal: "^^",
			},
		},
		{
			name:    "AND AND",
			literal: "AND",
			output: Token{
				Type:    AND,
				Literal: "AND",
			},
		},
		{
			name:    "AND and",
			literal: "and",
			output: Token{
				Type:    AND,
				Literal: "and",
			},
		},
		{
			name:    "OR OR",
			literal: "OR",
			output: Token{
				Type:    OR,
				Literal: "OR",
			},
		},
		{
			name:    "OR or",
			literal: "or",
			output: Token{
				Type:    OR,
				Literal: "or",
			},
		},
		{
			name:    "DATE 2018-01-01",
			literal: "2018-01-01",
			output: Token{
				Type:    DATE,
				Literal: "2018-01-01",
			},
		},
		{
			name:    "NUMBER 100",
			literal: "100",
			output: Token{
				Type:    NUMBER,
				Literal: "100",
			},
		},
		{
			name:    "NUMBER 1.0",
			literal: "1.0",
			output: Token{
				Type:    NUMBER,
				Literal: "1.0",
			},
		},
		{
			name:    "STRING a string",
			literal: "a string",
			output: Token{
				Type:    STRING,
				Literal: "a string",
			},
		},
		{
			name:    "EOF",
			literal: "",
			output: Token{
				Type:    EOF,
				Literal: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tok := New(test.literal)
			require.Equal(t, tok, test.output)
		})
	}
}
