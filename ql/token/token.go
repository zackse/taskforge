package token

import (
	"fmt"
	"strings"
)

// Type is used to determine which type of token this is
type Type int

// The supported token types
const (
	GT Type = iota
	LT
	GTE
	LTE
	EQ
	NE
	LIKE
	NLIKE

	AND
	OR

	LPAREN
	RPAREN

	EOF

	STRING
	NUMBER
	DATE

	ILLEGAL
	UNEXPECTED
)

func (t Type) String() string {
	switch t {
	case GT:
		return "GT"
	case LT:
		return "LT"
	case GTE:
		return "GTE"
	case LTE:
		return "LTE"
	case EQ:
		return "EQ"
	case NE:
		return "NE"
	case LIKE:
		return "LIKE"
	case NLIKE:
		return "NLIKE"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case EOF:
		return "EOF"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case DATE:
		return "DATE"
	case ILLEGAL:
		return "ILLEGAL"
	case UNEXPECTED:
		return "UNEXPECTED"
	default:
		return "UNKOWN"
	}
}

// LookupKeyword will return the type appropriate for the given string, if it is
// a keyword will return the token type for that keyword. Otherwise will return
// the token type STRING
func LookupKeyword(value string) Type {
	switch value {
	case "or":
		fallthrough
	case "OR":
		return OR
	case "and":
		fallthrough
	case "AND":
		return AND
	}

	return STRING
}

// DateOrNumber returns the appropriate token type for value if it is a date or
// a number.
func DateOrNumber(value string) Type {
	if strings.Contains(value, "-") {
		return DATE
	}

	return NUMBER
}

// Token is a lexical token of input
type Token struct {
	Type    Type
	Literal string
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%s, \"%s\")", t.Type, t.Literal)
}
