use super::tokens::Token;

#[derive(Debug, Clone, PartialEq)]
pub struct AST {
    pub expression: Expression,
}

#[derive(Debug, Clone, PartialEq)]
pub struct InfixExpression {
    pub left: Box<Expression>,
    pub operator: Token,
    pub right: Box<Expression>,
}

#[derive(Debug, Clone, PartialEq)]
pub enum Expression {
    Infix(Box<InfixExpression>),
    NumLiteral(Token),
    StrLiteral(Token),
    DateLiteral(Token),
}
