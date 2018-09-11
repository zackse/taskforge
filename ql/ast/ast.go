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

package ast

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chasinglogic/taskforge/ql/token"
)

// valid dateFormats in queries
var dateFormats = []string{
	"2006-01-02 03:04:05 PM",
	"2006-01-02 03:04 PM",
	"2006-01-02 03:04:05PM",
	"2006-01-02 03:04PM",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
	"2006-01-02",
}

// AST is the Abstract Syntax Tree for a query
type AST struct {
	Expression Expression
}

func (a AST) String() string {
	return a.Expression.String()
}

// New returns an AST for expression
func (a AST) New(exp Expression) AST {
	return AST{
		Expression: exp,
	}
}

// Node is implemented by all nodes in the tree
type Node interface {
	TokenLiteral() token.Token
	String() string
}

// Expression is represents an AST node that evaluates to a value
type Expression interface {
	Node
	expression()
}

// Literal is a literal value
type Literal interface {
	GetValue() interface{}
}

// NewLiteral returns a literal expression based on the given token
func NewLiteral(tok token.Token) (Expression, error) {
	var expression Expression

	switch tok.Type {
	case token.BOOLEAN:
		expression = BooleanLiteral{
			Token: tok,
			Value: strings.ToLower(tok.Literal) == "true",
		}
	case token.DATE:
		var value time.Time
		var err error
		for i := range dateFormats {
			value, err = time.ParseInLocation(dateFormats[i], tok.Literal, time.Local)
			if err == nil {
				break
			}
		}

		if value.IsZero() {
			return nil, err
		}

		expression = DateLiteral{
			Token: tok,
			Value: value,
		}
	case token.NUMBER:
		value, err := strconv.ParseFloat(tok.Literal, 64)
		if err != nil {
			return nil, err
		}

		expression = NumberLiteral{
			Token: tok,
			Value: value,
		}
	case token.STRING:
		expression = StringLiteral{
			Token: tok,
			Value: tok.Literal,
		}
	}

	return expression, nil
}

// InfixExpression is an infix AST node
type InfixExpression struct {
	Left     Expression
	Right    Expression
	Operator token.Token
}

func (ie InfixExpression) expression() {}

// TokenLiteral implements Node
func (ie InfixExpression) TokenLiteral() token.Token { return ie.Operator }

// String implements AST Node
func (ie InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator.Type.String() + " " + ie.Right.String() + ")"
}

// NumberLiteral is a literal number in a query
type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (nl NumberLiteral) expression() {}

// TokenLiteral implements Node
func (nl NumberLiteral) TokenLiteral() token.Token { return nl.Token }
func (nl NumberLiteral) String() string            { return fmt.Sprint(nl.Value) }

// GetValue returns the value for this literal
func (nl NumberLiteral) GetValue() interface{} { return nl.Value }

// BooleanLiteral is a literal boolean in a query
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (nl BooleanLiteral) expression() {}

// TokenLiteral implements Node
func (nl BooleanLiteral) TokenLiteral() token.Token { return nl.Token }
func (nl BooleanLiteral) String() string            { return fmt.Sprint(nl.Value) }

// GetValue returns the value for this literal
func (nl BooleanLiteral) GetValue() interface{} { return nl.Value }

// StringLiteral is a string
type StringLiteral struct {
	Token token.Token
	Value string
}

func (nl StringLiteral) expression() {}

// TokenLiteral implements Node
func (nl StringLiteral) TokenLiteral() token.Token { return nl.Token }
func (nl StringLiteral) String() string            { return "\"" + nl.Value + "\"" }

// GetValue returns the value for this literal
func (nl StringLiteral) GetValue() interface{} { return nl.Value }

// DateLiteral is a string
type DateLiteral struct {
	Token token.Token
	Value time.Time
}

func (nl DateLiteral) expression() {}

// TokenLiteral implements Node
func (nl DateLiteral) TokenLiteral() token.Token { return nl.Token }
func (nl DateLiteral) String() string            { return nl.Value.String() }

// GetValue returns the value for this literal
func (nl DateLiteral) GetValue() interface{} { return nl.Value }
