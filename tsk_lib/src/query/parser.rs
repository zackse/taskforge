use std::fmt;

use super::ast::{Expression, InfixExpression, AST};
use super::lexer::Lexer;
use super::tokens::Token;

type ExpResult = Result<Expression, ParseError>;

pub struct Validator;

impl Validator {
    fn string_field(field_name: &str, right: &Expression) -> Result<(), ParseError> {
        match right {
            Expression::StrLiteral(Token::Str(_)) => Ok(()),
            _ => Err(ParseError::from(format!(
                "{} can only be compared to a string got: {:?}",
                field_name, right
            ))),
        }
    }

    fn date_field(field_name: &str, right: &Expression) -> Result<(), ParseError> {
        match right {
            Expression::DateLiteral(Token::Date(_)) => Ok(()),
            _ => Err(ParseError::from(format!(
                "{} can only be compared to a string got: {:?}",
                field_name, right
            ))),
        }
    }

    fn number_field(field_name: &str, right: &Expression) -> Result<(), ParseError> {
        match right {
            Expression::NumLiteral(Token::Float(_)) => Ok(()),
            _ => Err(ParseError::from(format!(
                "{} can only be compared to a string got: {:?}",
                field_name, right
            ))),
        }
    }

    fn validate_comparison(infix: &InfixExpression) -> Result<(), ParseError> {
        match infix.left.as_ref() {
            Expression::StrLiteral(Token::Str(field)) => match field.as_ref() {
                "title" => Validator::string_field(&field, &infix.right),
                "context" => Validator::string_field(&field, &infix.right),
                "body" => Validator::string_field(&field, &infix.right),
                "notes" => Validator::string_field(&field, &infix.right),
                "created_date" => Validator::date_field(&field, &infix.right),
                "completed_date" => Validator::date_field(&field, &infix.right),
                "priority" => Validator::number_field(&field, &infix.right),
                _ => Err(ParseError::from(format!("invalid field name: {}", field))),
            },
            _ => Err(ParseError::from(format!(
                "invalid field expression: {:?}",
                infix.left
            ))),
        }
    }

    fn validate_logical(infix: &InfixExpression) -> Result<(), ParseError> {
        match infix.left.as_ref() {
            Expression::Infix(_) => Ok(()),
            _ => Err(ParseError::new(
                "logical operators can only compare other infix expressions",
            )),
        }
    }

    fn validate(infix: &InfixExpression) -> Result<(), ParseError> {
        match infix.operator {
            Token::AND | Token::OR => Validator::validate_logical(infix),
            _ => Validator::validate_comparison(infix),
        }
    }
}

pub struct Parser<'a> {
    lexer: Lexer<'a>,
    peek_token: Option<Token>,
}

#[derive(Debug)]
pub struct ParseError {
    pos: usize,
    ch: char,
    msg: String,
}

impl ParseError {
    pub fn new(msg: &str) -> ParseError {
        ParseError {
            pos: 0,
            ch: char::from(0),
            msg: msg.to_string(),
        }
    }

    pub fn at(mut self, pos: usize) -> ParseError {
        self.pos = pos;
        self
    }

    pub fn bad_char(mut self, ch: char) -> ParseError {
        self.ch = ch;
        self
    }
}

impl fmt::Display for ParseError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "ERROR: {} @ {}", self.msg, self.ch)
    }
}

impl From<String> for ParseError {
    fn from(input: String) -> ParseError {
        ParseError::new(&input)
    }
}

impl<'a> From<&'a str> for Parser<'a> {
    fn from(input: &'a str) -> Parser {
        let mut p = Parser {
            lexer: Lexer::from(input),
            peek_token: None,
        };

        // populate peek_token
        p.next();
        p
    }
}

impl<'a> Iterator for Parser<'a> {
    type Item = Token;

    fn next(&mut self) -> Option<Token> {
        let current_token = self.peek_token.clone();
        self.peek_token = self.lexer.next();
        println!("token: {:?}", current_token);
        current_token
    }
}

#[derive(PartialEq, PartialOrd)]
enum Precedence {
    LOWEST,
    STRING,
    ANDOR,
    COMPARISON,
}

impl<'a> From<&'a Token> for Precedence {
    fn from(token: &Token) -> Precedence {
        match token {
            Token::GT => Precedence::COMPARISON,
            Token::LT => Precedence::COMPARISON,
            Token::GTE => Precedence::COMPARISON,
            Token::LTE => Precedence::COMPARISON,
            Token::EQ => Precedence::COMPARISON,
            Token::NE => Precedence::COMPARISON,
            Token::AND => Precedence::ANDOR,
            Token::OR => Precedence::ANDOR,
            Token::Str(_) => Precedence::STRING,
            _ => Precedence::LOWEST,
        }
    }
}

impl<'a> From<&'a Option<Token>> for Precedence {
    fn from(token: &Option<Token>) -> Precedence {
        match token {
            Some(t) => Precedence::from(t),
            _ => Precedence::LOWEST,
        }
    }
}

impl<'a> Parser<'a> {
    pub fn parse(&mut self) -> Result<AST, ParseError> {
        self.parse_expression(Precedence::LOWEST)
            .map(|exp| AST { expression: exp })
    }

    fn parse_expression(&mut self, precedence: Precedence) -> ExpResult {
        println!("Parsing expression");

        let mut left = match self.next() {
            Some(token @ Token::Str(_)) => Expression::StrLiteral(token),
            Some(token @ Token::Date(_)) => Expression::DateLiteral(token),
            Some(token @ Token::Float(_)) => Expression::NumLiteral(token),
            Some(Token::LP) => self.parse_grouped_expression()?,
            Some(token) => return Err(ParseError::from(format!("no prefix func found: {}", token))),
            None => return Err(ParseError::new("unexpected end of input")),
        };

        println!("left: {:?}", left);

        while self.peek_token.is_some()
            && (precedence < Precedence::from(&self.peek_token) || precedence == Precedence::STRING)
        {
            println!("in loops");
            left = match self.peek_token {
                Some(Token::EQ) => self.parse_infix_exp(left),
                Some(Token::NE) => self.parse_infix_exp(left),
                Some(Token::LT) => self.parse_infix_exp(left),
                Some(Token::GT) => self.parse_infix_exp(left),
                Some(Token::LTE) => self.parse_infix_exp(left),
                Some(Token::GTE) => self.parse_infix_exp(left),
                Some(Token::AND) => self.parse_infix_exp(left),
                Some(Token::OR) => self.parse_infix_exp(left),
                Some(Token::Str(_)) => self.concat(left),
                _ => break,
            }?;

            println!("left: {:?}", left);
        }

        Ok(left)
    }

    fn concat(&mut self, left: Expression) -> ExpResult {
        println!("CONCATTING");
        match left {
            Expression::StrLiteral(Token::Str(mut s)) => {
                let next_char = match self.next() {
                    Some(Token::Str(val)) => val,
                    _ => return Err(ParseError::new("expected a string or field")),
                };

                s.push_str(" ");
                s.push_str(&next_char);

                Ok(Expression::StrLiteral(Token::Str(s)))
            }
            _ => Err(ParseError::new(
                "Expected an Expression::StrLiteral. If using a multi-word string in comparison make sure to quote it.",
            )),
        }
    }

    fn parse_infix_exp(&mut self, left: Expression) -> ExpResult {
        let op = match self.next() {
            Some(token @ Token::GT) => token,
            Some(token @ Token::LT) => token,
            Some(token @ Token::GTE) => token,
            Some(token @ Token::LTE) => token,
            Some(token @ Token::NE) => token,
            Some(token @ Token::EQ) => token,
            Some(token @ Token::AND) => token,
            Some(token @ Token::OR) => token,
            Some(token) => {
                return Err(ParseError::from(format!(
                    "{} is not a valid operator",
                    token
                )))
            }
            None => return Err(ParseError::new("Attempted infix found EOF")),
        };

        let precedence = Precedence::from(&op);
        let right = self.parse_expression(precedence)?;

        let exp = InfixExpression {
            left: Box::new(left),
            operator: op,
            right: Box::new(right),
        };

        // Validate the expression
        Validator::validate(&exp)?;

        Ok(Expression::Infix(Box::new(exp)))
    }

    fn parse_grouped_expression(&mut self) -> ExpResult {
        let exp = self.parse_expression(Precedence::LOWEST)?;

        match self.peek_token {
            Some(Token::RP) => {
                self.next();
                Ok(exp)
            }
            Some(_) => Err(ParseError::new("unclosed group expression: missing )")),
            None => Err(ParseError::new("unexpected EOF parsing group expression")),
        }
    }
}

#[cfg(test)]
pub mod test {
    use super::*;

    #[test]
    fn test_simple_exp() {
        let input = "this is a simple query";
        match Parser::from(input).parse() {
            Ok(ast) => assert_eq!(
                ast,
                AST {
                    expression: Expression::StrLiteral(Token::from("this is a simple query"))
                }
            ),
            Err(e) => {
                println!("{}", e);
                assert!(false)
            }
        }
    }

    #[test]
    fn test_simple_comparison() {
        let input = "title = something";
        match Parser::from(input).parse() {
            Ok(ast) => assert_eq!(
                ast,
                AST {
                    expression: Expression::Infix(Box::new(InfixExpression {
                        left: Box::new(Expression::StrLiteral(Token::from("title"))),
                        right: Box::new(Expression::StrLiteral(Token::from("something"))),
                        operator: Token::from('='),
                    }))
                }
            ),
            Err(e) => {
                println!("{}", e);
                assert!(false)
            }
        }
    }

    #[test]
    fn test_validator_invalid() {
        let exp = InfixExpression {
            left: Box::new(Expression::StrLiteral(Token::from("title"))),
            operator: Token::from('='),
            right: Box::new(Expression::NumLiteral(Token::from("1.0"))),
        };

        match Validator::validate(&exp) {
            Ok(()) => {
                println!("{:?}", exp);
                assert!(false)
            }
            Err(_) => (),
        }
    }

    #[test]
    fn test_invalid_comparison() {
        let input = "title = 1.0";
        match Parser::from(input).parse() {
            Ok(_) => assert!(false),
            Err(_) => (),
        }
    }

    #[test]
    fn test_logical_exp() {
        let input = "title = something and priority > 5";
        match Parser::from(input).parse() {
            Ok(ast) => {
                println!("AST: {:?}", ast);
                assert_eq!(
                    ast,
                    AST {
                        expression: Expression::Infix(Box::new(InfixExpression {
                            left: Box::new(Expression::Infix(Box::new(InfixExpression {
                                left: Box::new(Expression::StrLiteral(Token::from("title"))),
                                right: Box::new(Expression::StrLiteral(Token::from("something"))),
                                operator: Token::from('='),
                            }))),
                            operator: Token::AND,
                            right: Box::new(Expression::Infix(Box::new(InfixExpression {
                                left: Box::new(Expression::StrLiteral(Token::from("priority"))),
                                right: Box::new(Expression::NumLiteral(Token::Float(5.0))),
                                operator: Token::from('>'),
                            })))
                        }))
                    }
                )
            }
            Err(e) => {
                println!("{}", e);
                assert!(false)
            }
        }
    }

    #[test]
    fn test_complex_exp() {
        let input = "(title = something and priority > 5) or notes = \"what I want in notes\"";
        match Parser::from(input).parse() {
            Ok(ast) => {
                println!("AST: {:?}", ast);
                assert_eq!(
                    ast,
                    AST {
                        expression: Expression::Infix(Box::new(InfixExpression {
                            right: Box::new(Expression::Infix(Box::new(InfixExpression {
                                left: Box::new(Expression::StrLiteral(Token::from("notes"))),
                                right: Box::new(Expression::StrLiteral(Token::from(
                                    "what I want in notes"
                                ))),
                                operator: Token::EQ,
                            }))),
                            operator: Token::OR,
                            left: Box::new(Expression::Infix(Box::new(InfixExpression {
                                left: Box::new(Expression::Infix(Box::new(InfixExpression {
                                    left: Box::new(Expression::StrLiteral(Token::from("title"))),
                                    right: Box::new(Expression::StrLiteral(Token::from(
                                        "something"
                                    ))),
                                    operator: Token::from('='),
                                }))),
                                operator: Token::AND,
                                right: Box::new(Expression::Infix(Box::new(InfixExpression {
                                    left: Box::new(Expression::StrLiteral(Token::from("priority"))),
                                    right: Box::new(Expression::NumLiteral(Token::Float(5.0))),
                                    operator: Token::from('>'),
                                })))
                            })))
                        }))
                    }
                )
            }
            Err(e) => {
                println!("{}", e);
                assert!(false)
            }
        }
    }
}
