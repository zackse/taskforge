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
        if let Some(next_char) = self.content.peek() {
            if end(*next_char) {
                return s;
            }
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
            (c.is_alphabetic() && c != '.' && c != '-') || c.is_whitespace() || c == ')'
        });

        if s.contains("-") {
            Token::from(s.as_ref())
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

    fn next_char(&mut self) -> Option<char> {
        self.current_char = self.content.next();
        self.current_char
    }

    fn token_from_char(&mut self, c: char) -> Token {
        match c {
            '>' | '<' => if let Some('=') = self.content.peek() {
                let mut s = c.to_string();
                // Advance past the = sign.
                // Safe to unwrap since the peek already showed us it was Some
                s.push(self.next_char().unwrap());
                Token::from(s.as_ref())
            } else {
                Token::from(c)
            },
            '"' => self.string(),
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

    #[test]
    fn test_all_unquoted_strings() {
        let input = "this is a simple query";
        let tokens: Vec<Token> = Lexer::from(input).collect();
        assert_eq!(tokens[0], Token::Str("this".to_string()));
        assert_eq!(tokens[1], Token::Str("is".to_string()));
        assert_eq!(tokens[2], Token::Str("a".to_string()));
        assert_eq!(tokens[3], Token::Str("simple".to_string()));
        assert_eq!(tokens[4], Token::Str("query".to_string()));
        assert_eq!(tokens.len(), 5);
    }

    #[test]
    fn test_unquoted_string_op_string() {
        let input = "unquoted_string = \"this is a string\"";
        let tokens: Vec<Token> = Lexer::from(input).collect();
        assert_eq!(tokens[0], Token::Str("unquoted_string".to_string()));
        assert_eq!(tokens[1], Token::EQ);
        assert_eq!(tokens[2], Token::Str("this is a string".to_string()));
        assert_eq!(tokens.len(), 3);
    }

    #[test]
    fn test_unquoted_string_complex_op_string() {
        let input = "unquoted_string >= \"this is a string\"";
        let tokens: Vec<Token> = Lexer::from(input).collect();
        assert_eq!(tokens[0], Token::Str("unquoted_string".to_string()));
        assert_eq!(tokens[1], Token::GTE);
        assert_eq!(tokens[2], Token::Str("this is a string".to_string()));
        assert_eq!(tokens.len(), 3);
    }

    #[test]
    fn test_num_lexing() {
        let mut input = "5";
        let mut tokens: Vec<Token> = Lexer::from(input).collect();
        assert_eq!(tokens[0], Token::Float(5.0));

        input = "(priority > 5)";
        tokens = Lexer::from(input).collect();
        assert_eq!(tokens[0], Token::LP);
        assert_eq!(tokens[1], Token::from("priority"));
        assert_eq!(tokens[2], Token::from('>'));
        assert_eq!(tokens[3], Token::Float(5.0));
        assert_eq!(tokens[4], Token::RP);

        input = "(priority > 5.0)";
        tokens = Lexer::from(input).collect();
        assert_eq!(tokens[0], Token::LP);
        assert_eq!(tokens[1], Token::from("priority"));
        assert_eq!(tokens[2], Token::from('>'));
        assert_eq!(tokens[3], Token::Float(5.0));
        assert_eq!(tokens[4], Token::RP);
    }

    #[test]
    fn test_complex_exp() {
        let input = "(unquoted_string >= 5.0 and other = \"other string\") or (mighty morphin power rangers)";
        let tokens: Vec<Token> = Lexer::from(input).collect();
        assert_eq!(tokens[0], Token::LP);
        assert_eq!(tokens[1], Token::Str("unquoted_string".to_string()));
        assert_eq!(tokens[2], Token::GTE);
        assert_eq!(tokens[3], Token::Float(5.0));
        assert_eq!(tokens[4], Token::AND);
        assert_eq!(tokens[5], Token::Str("other".to_string()));
        assert_eq!(tokens[6], Token::EQ);
        assert_eq!(tokens[7], Token::Str("other string".to_string()));
        assert_eq!(tokens[8], Token::RP);
        assert_eq!(tokens[9], Token::OR);
        assert_eq!(tokens[10], Token::LP);
        assert_eq!(tokens[11], Token::Str("mighty".to_string()));
        assert_eq!(tokens[12], Token::Str("morphin".to_string()));
        assert_eq!(tokens[13], Token::Str("power".to_string()));
        assert_eq!(tokens[14], Token::Str("rangers".to_string()));
        assert_eq!(tokens[15], Token::RP);
        assert_eq!(tokens.len(), 16);
    }
}
