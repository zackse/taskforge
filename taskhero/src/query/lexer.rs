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

        loop {
            match self.content.next() {
                Some(c) => {
                    if end(c) {
                        break;
                    }

                    s.push(c);
                }
                None => break,
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
        let s = self.read_until(|c| (c.is_alphabetic() && c != '.') || c.is_whitespace());
        match s.parse::<f64>() {
            Ok(flt) => Token::Float(flt),
            Err(e) => Token::Unexpected(format!("{}", e)),
        }
    }

    fn field(&mut self) -> Token {
        let s = self.read_until(|c| !c.is_alphabetic());
        Token::from(s.as_ref())
    }

    fn next_char(&mut self) -> Option<char> {
        self.current_char = self.content.next();
        self.current_char
    }
}

impl<'a> Iterator for Lexer<'a> {
    type Item = Token;

    fn next(&mut self) -> Option<Token> {
        match self.next_char() {
            Some(c) => Some(match c {
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
                'a'..='z' | 'A'..='Z' => self.field(),
                _num if c.is_digit(10) => self.number(),
                _whitespace if c.is_whitespace() => return self.next(),
                _ => Token::from(c),
            }),
            None => None,
        }
    }
}

#[cfg(test)]
pub mod tests {
    use super::*;

    #[test]
    fn test_all_fields() {
        let input = "some more fields";
        let tokens: Vec<Token> = Lexer::from(input).collect();
        assert_eq!(tokens[0], Token::Field("some".to_string()));
        assert_eq!(tokens[1], Token::Field("more".to_string()));
        assert_eq!(tokens[2], Token::Field("fields".to_string()));
        assert_eq!(tokens.len(), 3);
    }

    #[test]
    fn test_field_op_string() {
        let input = "field = \"this is a string\"";
        let tokens: Vec<Token> = Lexer::from(input).collect();
        assert_eq!(tokens[0], Token::Field("field".to_string()));
        assert_eq!(tokens[1], Token::EQ);
        assert_eq!(tokens[2], Token::Str("this is a string".to_string()));
        assert_eq!(tokens.len(), 3);
    }

    #[test]
    fn test_field_complex_op_string() {
        let input = "field >= \"this is a string\"";
        let tokens: Vec<Token> = Lexer::from(input).collect();
        assert_eq!(tokens[0], Token::Field("field".to_string()));
        assert_eq!(tokens[1], Token::GTE);
        assert_eq!(tokens[2], Token::Str("this is a string".to_string()));
        assert_eq!(tokens.len(), 3);
    }

    #[test]
    fn test_complex_exp() {
        let input = "(field >= 5.0 and other = \"other string\") or (mighty morphin power rangers)";
        let tokens: Vec<Token> = Lexer::from(input).collect();
        assert_eq!(tokens[0], Token::LP);
        assert_eq!(tokens[1], Token::Field("field".to_string()));
        assert_eq!(tokens[2], Token::GTE);
        assert_eq!(tokens[3], Token::Float(5.0));
        assert_eq!(tokens[4], Token::AND);
        assert_eq!(tokens[5], Token::Field("other".to_string()));
        assert_eq!(tokens[6], Token::EQ);
        assert_eq!(tokens[7], Token::Str("other string".to_string()));
        assert_eq!(tokens.len(), 16);
    }
}
