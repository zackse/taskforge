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


// Package lexer contains the lexer for our query language
package lexer

import (
	"github.com/chasinglogic/tsk/ql/token"
)

// Lexer maintains document position and lexing the input
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

// New create a new lexer for the given input
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Char returns the current character under cursor
func (l *Lexer) Char() byte {
	return l.ch
}

// Pos returns the current lexer position
func (l *Lexer) Pos() int {
	return l.position
}

func (l *Lexer) String() string {
	if len(l.input) <= l.position {
		return l.input
	}

	output := "\n" + l.input + "\n"

	for i := 0; i < l.position; i++ {
		output += " "
	}

	output += "^"
	return output
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) read(valid func(byte) bool) string {
	start := l.position

	for valid(l.ch) {
		l.readChar()
	}

	return l.input[start:l.position]
}

func (l *Lexer) number() token.Token {
	tok := token.Token{
		Literal: l.read(func(ch byte) bool {
			return isDigit(ch) || ch == '-' || ch == '.' || ch == ':'
		}),
	}

	tok.Type = token.DateOrNumber(tok.Literal)

	return tok
}

func (l *Lexer) unquotedString() token.Token {
	tok := token.Token{
		Type: token.STRING,
		Literal: l.read(func(ch byte) bool {
			return isLetter(ch) || isDigit(ch)
		}),
	}

	return tok
}

func (l *Lexer) quotedString() token.Token {
	// Skip opening quote
	l.readChar()

	tok := token.Token{
		Type: token.STRING,
		Literal: l.read(func(ch byte) bool {
			return isLetter(ch) || ch == ' ' || isDigit(ch)
		}),
	}

	// Check for closing quote
	if l.ch != '"' {
		return token.Token{Type: token.UNEXPECTED, Literal: string(l.ch)}
	}

	return tok
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

// NextToken returns the next token.Token found in the given input
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch {
	case l.ch == '=':
		tok = token.Token{Type: token.EQ, Literal: string(l.ch)}
	case l.ch == '<' || l.ch == '>':
		tok = token.Token{Literal: string(l.ch)}
		if l.ch == '<' {
			tok.Type = token.LT
		} else {
			tok.Type = token.GT
		}

		if l.peekChar() == '=' {
			l.readChar()
			tok.Literal += "="
			if tok.Type == token.LT {
				tok.Type = token.LTE
			} else {
				tok.Type = token.GTE
			}
		}
	case l.ch == '!' || l.ch == '^':
		switch l.peekChar() {
		case '=':
			l.readChar()
			tok = token.Token{Type: token.NE, Literal: "!="}
		case '~':
			l.readChar()
			tok = token.Token{Type: token.NLIKE, Literal: "!~"}
		case ' ':
			tok = token.Token{Type: token.LIKE, Literal: "^"}
		default:
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}
		}
	case l.ch == '~':
		tok = token.Token{Type: token.LIKE, Literal: string(l.ch)}
	case l.ch == '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch)}
	case l.ch == ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch)}
	case l.ch == '"':
		tok = l.quotedString()
	case '0' <= l.ch && l.ch <= '9':
		tok = l.number()
	case l.ch == '-':
		// skip the -
		l.readChar()
		tok = l.unquotedString()
	case 'a' <= l.ch && l.ch <= 'z' || 'A' <= l.ch && l.ch <= 'Z':
		tok = l.unquotedString()
		tok.Type = token.LookupKeyword(tok.Literal)
	case l.ch == 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}
	}

	l.readChar()
	return tok
}

func isWhitespace(ch byte) bool {
	return ' ' == ch || '\n' == ch || '\t' == ch
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '-' || ch == ','
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
