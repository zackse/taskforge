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
	BOOLEAN

	ILLEGAL
	UNEXPECTED
)

func (t Type) String() string {
	switch t {
	case GT:
		return ">"
	case LT:
		return "<"
	case GTE:
		return ">="
	case LTE:
		return "<="
	case EQ:
		return "="
	case NE:
		return "!="
	case LIKE:
		return "~"
	case NLIKE:
		return "!~"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case EOF:
		return "EOF"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case DATE:
		return "DATE"
	case BOOLEAN:
		return "BOOLEAN"
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
	case "false":
		fallthrough
	case "False":
		fallthrough
	case "true":
		fallthrough
	case "True":
		return BOOLEAN
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
