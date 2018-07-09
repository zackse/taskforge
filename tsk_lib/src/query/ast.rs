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
