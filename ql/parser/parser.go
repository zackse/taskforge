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
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/chasinglogic/tsk/ql/ast"
	"github.com/chasinglogic/tsk/ql/lexer"
	"github.com/chasinglogic/tsk/ql/token"
)

// Operator order of precedence
const (
	_      int = iota
	LOWEST     // Lowest priority
	STRING
	ANDOR
	COMPARISON
)

var precedences = map[token.Type]int{
	token.EQ:     COMPARISON,
	token.NE:     COMPARISON,
	token.LT:     COMPARISON,
	token.GT:     COMPARISON,
	token.GTE:    COMPARISON,
	token.LTE:    COMPARISON,
	token.LIKE:   COMPARISON,
	token.NLIKE:  COMPARISON,
	token.AND:    ANDOR,
	token.OR:     ANDOR,
	token.STRING: STRING,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser parses a PQL query and returns an AST for that query
type Parser struct {
	l      *lexer.Lexer
	errors []error

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

// New returns a parser for the given lexer
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []error{},
	}

	p.prefixParseFns = map[token.Type]prefixParseFn{
		token.STRING: p.parseString,
		token.NUMBER: p.parseNumberLiteral,
		token.DATE:   p.parseDateLiteral,
		token.LPAREN: p.parseGroupedExpression,
	}

	p.infixParseFns = map[token.Type]infixParseFn{
		token.EQ:     p.parseInfixExpression,
		token.NE:     p.parseInfixExpression,
		token.LT:     p.parseInfixExpression,
		token.GT:     p.parseInfixExpression,
		token.GTE:    p.parseInfixExpression,
		token.LTE:    p.parseInfixExpression,
		token.LIKE:   p.parseInfixExpression,
		token.NLIKE:  p.parseInfixExpression,
		token.OR:     p.parseLogicExpression,
		token.AND:    p.parseLogicExpression,
		token.STRING: p.concat,
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

// Error will return a parser.Error instance with the problems found during
// parsing
func (p *Parser) Error() error {
	if len(p.errors) != 0 {
		return fmt.Errorf("parsing errors: %v", p.errors)
	}

	return nil
}

func (p *Parser) addError(err error) {
	p.errors = append(p.errors,
		fmt.Errorf("::%s:%d:%s", string(p.l.Char()), p.l.Pos(), err.Error()))
}

func (p *Parser) peekError(t token.Type) {
	p.addError(fmt.Errorf("expected next token to be %s, got %s instead", t, p.peekToken.Type))
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	p.addError(fmt.Errorf("%s not allowed in comparison expression", t))
}

// Parse will turn the given query into an ast.AST
// TODO: Write this
func (p *Parser) Parse() ast.AST {
	var a ast.AST
	a.Expression = p.parseExpression(LOWEST)
	return a
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if p.curToken.Type == token.EOF {
		p.addError(errors.New("unexpected end of input"))
		return nil
	}

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.EOF) &&
		precedence < p.peekPrecedence() {

		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseString() ast.Expression {
	return ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) concat(left ast.Expression) ast.Expression {
	sl, ok := left.(ast.StringLiteral)
	if !ok {
		return left
	}

	sl.Token.Literal += " " + p.curToken.Literal
	sl.Value += " " + p.curToken.Literal
	return sl
}

var dateFormats = []string{
	"2006-01-02 03:04:05 PM",
	"2006-01-02 03:04 PM",
	"2006-01-02 03:04:05PM",
	"2006-01-02 03:04PM",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
	"2006-01-02",
}

func (p *Parser) parseDateLiteral() ast.Expression {
	lit := ast.DateLiteral{Token: p.curToken}

	for i := range dateFormats {
		value, err := time.Parse(dateFormats[i], p.curToken.Literal)
		if err == nil {
			lit.Value = value
			break
		}
	}

	if lit.Value.IsZero() {
		p.addError(fmt.Errorf("not a valid date format: %s", p.curToken.Literal))
	}

	return lit
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	lit := ast.NumberLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addError(fmt.Errorf("could not parse float: %s", p.curToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseLogicExpression(left ast.Expression) ast.Expression {
	expression := ast.InfixExpression{
		Operator: p.curToken,
		Left:     left,
	}

	switch expression.Left.(type) {
	case ast.StringLiteral:
		break
	case ast.InfixExpression:
		break
	default:
		p.addError(errors.New("logic operators must (AND / OR) must be followed by a comparison or string literal expression"))
		return nil
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	switch expression.Right.(type) {
	case ast.StringLiteral:
		break
	case ast.InfixExpression:
		break
	default:
		p.addError(errors.New("logic operators must (AND / OR) must be followed by a comparison or string literal expression"))
		return nil
	}

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := ast.InfixExpression{
		Operator: p.curToken,
		Left:     left,
	}

	if _, ok := left.(ast.StringLiteral); !ok {
		p.addError(fmt.Errorf("left side of an infix expression must be a string literal got %s", left))
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}
