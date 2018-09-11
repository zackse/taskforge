package ast

import (
	"testing"
	"time"

	"github.com/chasinglogic/taskforge/ql/token"
	"github.com/stretchr/testify/require"
)

func TestNewBooleanLiteral(t *testing.T) {
	tests := []struct {
		name     string
		token    token.Token
		expected Expression
	}{
		{
			name: "boolean literal true",
			token: token.Token{
				Type:    token.BOOLEAN,
				Literal: "true",
			},
			expected: BooleanLiteral{
				Token: token.Token{
					Type:    token.BOOLEAN,
					Literal: "true",
				},
				Value: true,
			},
		},
		{
			name: "boolean literal True",
			token: token.Token{
				Type:    token.BOOLEAN,
				Literal: "True",
			},
			expected: BooleanLiteral{
				Token: token.Token{
					Type:    token.BOOLEAN,
					Literal: "True",
				},
				Value: true,
			},
		},
		{
			name: "boolean literal False",
			token: token.Token{
				Type:    token.BOOLEAN,
				Literal: "False",
			},
			expected: BooleanLiteral{
				Token: token.Token{
					Type:    token.BOOLEAN,
					Literal: "False",
				},
				Value: false,
			},
		},
		{
			name: "boolean literal false",
			token: token.Token{
				Type:    token.BOOLEAN,
				Literal: "false",
			},
			expected: BooleanLiteral{
				Token: token.Token{
					Type:    token.BOOLEAN,
					Literal: "false",
				},
				Value: false,
			},
		},
		{
			name: "number literal 1",
			token: token.Token{
				Type:    token.NUMBER,
				Literal: "1",
			},
			expected: NumberLiteral{
				Token: token.Token{
					Type:    token.NUMBER,
					Literal: "1",
				},
				Value: 1.0,
			},
		},
		{
			name: "number literal 1",
			token: token.Token{
				Type:    token.NUMBER,
				Literal: "1",
			},
			expected: NumberLiteral{
				Token: token.Token{
					Type:    token.NUMBER,
					Literal: "1",
				},
				Value: 1.0,
			},
		},
		{
			name: "date literal format 2006-01-02 03:04:05 PM",
			token: token.Token{
				Type:    token.DATE,
				Literal: "2018-01-01 01:01:01 PM",
			},
			expected: DateLiteral{
				Token: token.Token{
					Type:    token.DATE,
					Literal: "2018-01-01 01:01:01 PM",
				},
				Value: time.Date(2018, time.January, 1, 13, 1, 1, 0, time.Local),
			},
		},
		{
			name: "date literal format 2006-01-02 03:04:05PM",
			token: token.Token{
				Type:    token.DATE,
				Literal: "2018-01-01 01:01:01PM",
			},
			expected: DateLiteral{
				Token: token.Token{
					Type:    token.DATE,
					Literal: "2018-01-01 01:01:01PM",
				},
				Value: time.Date(2018, time.January, 1, 13, 1, 1, 0, time.Local),
			},
		},
		{
			name: "date literal format 2006-01-02 03:04 PM",
			token: token.Token{
				Type:    token.DATE,
				Literal: "2018-01-01 01:01 PM",
			},
			expected: DateLiteral{
				Token: token.Token{
					Type:    token.DATE,
					Literal: "2018-01-01 01:01 PM",
				},
				Value: time.Date(2018, time.January, 1, 13, 1, 0, 0, time.Local),
			},
		},
		{
			name: "date literal format 2006-01-02 03:04PM",
			token: token.Token{
				Type:    token.DATE,
				Literal: "2018-01-01 01:01PM",
			},
			expected: DateLiteral{
				Token: token.Token{
					Type:    token.DATE,
					Literal: "2018-01-01 01:01PM",
				},
				Value: time.Date(2018, time.January, 1, 13, 1, 0, 0, time.Local),
			},
		},

		{
			name: "date literal format 2006-01-02 03:04:05 AM",
			token: token.Token{
				Type:    token.DATE,
				Literal: "2018-01-01 01:01:01 AM",
			},
			expected: DateLiteral{
				Token: token.Token{
					Type:    token.DATE,
					Literal: "2018-01-01 01:01:01 AM",
				},
				Value: time.Date(2018, time.January, 1, 1, 1, 1, 0, time.Local),
			},
		},
		{
			name: "date literal format 2006-01-02 03:04:05AM",
			token: token.Token{
				Type:    token.DATE,
				Literal: "2018-01-01 01:01:01AM",
			},
			expected: DateLiteral{
				Token: token.Token{
					Type:    token.DATE,
					Literal: "2018-01-01 01:01:01AM",
				},
				Value: time.Date(2018, time.January, 1, 1, 1, 1, 0, time.Local),
			},
		},
		{
			name: "date literal format 2006-01-02 03:04 AM",
			token: token.Token{
				Type:    token.DATE,
				Literal: "2018-01-01 01:01 AM",
			},
			expected: DateLiteral{
				Token: token.Token{
					Type:    token.DATE,
					Literal: "2018-01-01 01:01 AM",
				},
				Value: time.Date(2018, time.January, 1, 1, 1, 0, 0, time.Local),
			},
		},
		{
			name: "date literal format 2006-01-02 03:04AM",
			token: token.Token{
				Type:    token.DATE,
				Literal: "2018-01-01 01:01AM",
			},
			expected: DateLiteral{
				Token: token.Token{
					Type:    token.DATE,
					Literal: "2018-01-01 01:01AM",
				},
				Value: time.Date(2018, time.January, 1, 1, 1, 0, 0, time.Local),
			},
		},
		{
			name: "date literal format 2006-01-02 15:04:05",
			token: token.Token{
				Type:    token.DATE,
				Literal: "2018-01-01 01:01:01",
			},
			expected: DateLiteral{
				Token: token.Token{
					Type:    token.DATE,
					Literal: "2018-01-01 01:01:01",
				},
				Value: time.Date(2018, time.January, 1, 1, 1, 1, 0, time.Local),
			},
		},
		{
			name: "date literal format 2006-01-02 15:04",
			token: token.Token{
				Type:    token.DATE,
				Literal: "2018-01-01 01:01",
			},
			expected: DateLiteral{
				Token: token.Token{
					Type:    token.DATE,
					Literal: "2018-01-01 01:01",
				},
				Value: time.Date(2018, time.January, 1, 1, 1, 0, 0, time.Local),
			},
		},
		{
			name: "date literal format 2006-01-02",
			token: token.Token{
				Type:    token.DATE,
				Literal: "2018-01-01",
			},
			expected: DateLiteral{
				Token: token.Token{
					Type:    token.DATE,
					Literal: "2018-01-01",
				},
				Value: time.Date(2018, time.January, 1, 0, 0, 0, 0, time.Local),
			},
		},
		{
			name: "string literal string literal",
			token: token.Token{
				Type:    token.STRING,
				Literal: "string literal",
			},
			expected: StringLiteral{
				Token: token.Token{
					Type:    token.STRING,
					Literal: "string literal",
				},
				Value: "string literal",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			literal, err := NewLiteral(test.token)
			require.NoError(t, err)
			require.Equal(t, literal, test.expected)
		})
	}
}

func TestNewLiteralNilIfInvalidTokenType(t *testing.T) {
	exp, _ := NewLiteral(token.Token{
		Type:    token.LPAREN,
		Literal: "(",
	})

	require.Nil(t, exp)
}
