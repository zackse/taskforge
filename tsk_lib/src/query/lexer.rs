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

use super::tokens::Token;
use std::iter;
use std::str;

pub struct Lexer<'a> {
    content: iter::Peekable<str::Chars<'a>>,
    current_char: Option<char>,
}

impl<'a> From<&'a str> for Lexer<'a> {
    fn from(input: &'a str) -> Lexer {
        Lexer {
            content: input.chars().peekable(),
            current_char: None,
        }
    }
}

impl<'a> Lexer<'a> {
    fn read_until<T>(&mut self, end: T) -> String
    where
        T: Fn(char) -> bool,
    {
        // Safe to unwrap here since we already validated that current_char is not None as part of
        // next
        let mut s = self.current_char.unwrap().to_string();
        // Check for early end on one character tokens
        match self.content.peek() {
            Some(c) => if end(*c) {
                return s;
            },
            None => return s,
        }

        while let Some(c) = self.content.next() {
            s.push(c);

            // We have to check the peek character since we don't want to consume it
            if let Some(next_char) = self.content.peek() {
                if end(*next_char) {
                    break;
                }
            }
        }

        s
    }

    fn string(&mut self) -> Token {
        let mut s = self.read_until(|c| c == '"');
        s.remove(0);
        // Advance past the " char
        self.next_char();
        Token::Str(s)
    }

    fn number(&mut self) -> Token {
        let s = self.read_until(|c| {
            (c.is_alphabetic() && c != '.' && c != '-' && c != ':') || c.is_whitespace() || c == ')'
        });

        // It's a date if it has - in it
        if s.contains("-") {
            Token::from(s.as_ref())
        // Otherwise we need to remove whitespace for the float to properly parsed
        } else {
            Token::from(
                s.chars()
                    .filter(|c| !c.is_whitespace())
                    .collect::<String>()
                    .as_ref(),
            )
        }
    }

    fn unquoted_string(&mut self) -> Token {
        let s = self.read_until(|c| c.is_whitespace() || c == ')' || c == '(');
        Token::from(s.as_ref())
    }

    fn keyword_string(&mut self) -> Token {
        // Skip -
        self.next_char();

        let s = self.read_until(|c| c.is_whitespace() || c == ')' || c == '(');
        Token::Str(s)
    }

    fn next_char(&mut self) -> Option<char> {
        self.current_char = self.content.next();
        self.current_char
    }

    fn token_from_char(&mut self, c: char) -> Token {
        match c {
            '>' | '<' | '^' | '!' => match self.content.peek() {
                Some('=') => {
                    let mut s = c.to_string();
                    // Advance past the = sign.
                    // Safe to unwrap since the peek already showed us it was Some
                    s.push(self.next_char().unwrap());
                    Token::from(s.as_ref())
                }
                Some('^') if c == '^' => {
                    // Skip second ^
                    self.next_char();
                    Token::NLIKE
                }
                _ => Token::from(c),
            },
            '"' => self.string(),
            '-' => match self.content.peek() {
                Some('a'..='z') | Some('A'..='Z') => self.keyword_string(),
                _ => self.number(),
            },
            'a'..='z' | 'A'..='Z' => self.unquoted_string(),
            _num if c.is_digit(10) => self.number(),
            _ => Token::from(c),
        }
    }
}

impl<'a> Iterator for Lexer<'a> {
    type Item = Token;

    fn next(&mut self) -> Option<Token> {
        while let Some(c) = self.next_char() {
            if c.is_whitespace() {
                continue;
            }

            return Some(self.token_from_char(c));
        }

        None
    }
}

#[cfg(test)]
pub mod tests {
    use super::*;

    struct LexerTest<'a> {
        input: &'a str,
        name: &'a str,
        expected: Vec<Token>,
    }

    #[test]
    fn test_lexer() {
        let tests = vec![
            LexerTest {
                input: "this is a simple query",
                name: "all unquoted strings",
                expected: vec![
                    Token::Str("this".to_string()),
                    Token::Str("is".to_string()),
                    Token::Str("a".to_string()),
                    Token::Str("simple".to_string()),
                    Token::Str("query".to_string()),
                ],
            },
            LexerTest {
                input: "unquoted_string = \"this is a string\"",
                name: "unquoted string op quoted string",
                expected: vec![
                    Token::Str("unquoted_string".to_string()),
                    Token::EQ,
                    Token::Str("this is a string".to_string()),
                ],
            },
            LexerTest {
                input: "5",
                name: "one number",
                expected: vec![Token::Float(5.0)],
            },
            LexerTest {
                input: "(priority > 5)",
                name: "grouped priority greater than 5",
                expected: vec![
                    Token::LP,
                    Token::from("priority"),
                    Token::from('>'),
                    Token::Float(5.0),
                    Token::RP,
                ],
            },
            LexerTest {
                input: "(priority > 5.0)",
                name: "grouped priority greater than 5.0",
                expected: vec![
                    Token::LP,
                    Token::from("priority"),
                    Token::from('>'),
                    Token::Float(5.0),
                    Token::RP,
                ],
            },
            LexerTest {
                input:  "(unquoted_string >= 5.0 and other = \"other string\") or (mighty morphin power rangers)",
                name: "complex expression",
                expected: vec![
                    Token::LP,
                    Token::Str("unquoted_string".to_string()),
                    Token::GTE,
                    Token::Float(5.0),
                    Token::AND,
                    Token::Str("other".to_string()),
                    Token::EQ,
                    Token::Str("other string".to_string()),
                    Token::RP,
                    Token::OR,
                    Token::LP,
                    Token::Str("mighty".to_string()),
                    Token::Str("morphin".to_string()),
                    Token::Str("power".to_string()),
                    Token::Str("rangers".to_string()),
                    Token::RP,
                ],
            },
            LexerTest {
                input: "milk -and cookies",
                name: "escaped keyword",
                expected: vec![
                    Token::Str("milk".to_string()),
                    Token::Str("and".to_string()),
                    Token::Str("cookies".to_string()),
                ],
            },
            LexerTest {
                input: "title ^ \"take out the trash\"",
                name: "title like take out the trash",
                expected: vec![
                    Token::Str("title".to_string()),
                    Token::LIKE,
                    Token::Str("take out the trash".to_string()),
                ]
            },
            LexerTest {
                input: "(\"milk and sugar\") and priority > 5",
                name: "grouped string query and priority greater than 5",
                expected: vec![
                    Token::LP,
                    Token::Str("milk and sugar".to_string()),
                    Token::RP,
                    Token::AND,
                    Token::Str("priority".to_string()),
                    Token::GT,
                    Token::Float(5.0),
                ]
            },
            LexerTest {
                input: "(priority > 5 and title ^ \"take out the trash\") or (context = \"work\" and (priority >= 2 or (\"my little pony\")))",
                name: "very complex expression",
                expected: vec![
                    Token::LP,
                    Token::Str("priority".to_string()),
                    Token::GT,
                    Token::Float(5.0),
                    Token::AND,
                    Token::Str("title".to_string()),
                    Token::LIKE,
                    Token::Str("take out the trash".to_string()),
                    Token::RP,
                    Token::OR,
                    Token::LP,
                    Token::Str("context".to_string()),
                    Token::EQ,
                    Token::Str("work".to_string()),
                    Token::AND,
                    Token::LP,
                    Token::Str("priority".to_string()),
                    Token::GTE,
                    Token::Float(2.0),
                    Token::OR,
                    Token::LP,
                    Token::Str("my little pony".to_string()),
                    Token::RP,
                    Token::RP,
                    Token::RP
                ]
            }
        ];

        for test in tests.iter() {
            println!("Running Test: {}", test.name);

            let tokens: Vec<Token> = Lexer::from(test.input).collect();

            for ind in 0..test.expected.len() {
                let lexed_token = &tokens[ind];
                let expected_token = &test.expected[ind];
                assert_eq!(lexed_token, expected_token);
            }

            assert_eq!(tokens.len(), test.expected.len(),);
        }
    }
}
