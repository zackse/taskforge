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
	"github.com/chasinglogic/taskforge/ql/token"
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

func (l *Lexer) number() string {
	return l.read(func(ch byte) bool {
		return isDigit(ch) || ch == '-' || ch == '.' || ch == ':'
	})
}

func (l *Lexer) unquotedString() string {
	return l.read(func(ch byte) bool {
		return isLetter(ch) || isDigit(ch)
	})
}

func (l *Lexer) quotedString() string {
	// Skip opening quote
	l.readChar()

	literal := l.read(func(ch byte) bool {
		return isLetter(ch) || ch == ' ' || isDigit(ch)
	})

	// Check for closing quote
	if l.ch != '"' {
		return ""
	}

	return literal
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
	case l.ch == '<' || l.ch == '>' || l.ch == '!':
		literal := string(l.ch)
		if l.peekChar() == '=' {
			l.readChar()
			literal += "="
		}

		tok = token.New(literal)
	case l.ch == '^':
		switch l.peekChar() {
		case '=':
			l.readChar()
			tok = token.New(string(l.ch) + "=")
		case '^':
			l.readChar()
			tok = token.New(string(l.ch) + "^")
		default:
			tok = token.New(string(l.ch))
		}
	case l.ch == '"':
		tok = token.New(l.quotedString())
	case '0' <= l.ch && l.ch <= '9':
		return token.New(l.number())
	case l.ch == '-':
		// skip the -
		l.readChar()
		tok = token.New(l.unquotedString())
		tok.Type = token.STRING
	case 'a' <= l.ch && l.ch <= 'z' || 'A' <= l.ch && l.ch <= 'Z':
		tok = token.New(l.unquotedString())
	case l.ch == 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		tok = token.New(string(l.ch))
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
