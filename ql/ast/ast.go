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
	"time"

	"github.com/chasinglogic/tsk/ql/token"
)

// AST is the Abstract Syntax Tree for a query
type AST struct {
	Expression Expression
}

func (a AST) String() string {
	return a.Expression.String()
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

// StringLiteral is a string
type StringLiteral struct {
	Token token.Token
	Value string
}

func (nl StringLiteral) expression() {}

// TokenLiteral implements Node
func (nl StringLiteral) TokenLiteral() token.Token { return nl.Token }
func (nl StringLiteral) String() string            { return nl.Value }

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
