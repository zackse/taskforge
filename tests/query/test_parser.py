# pylint: disable=missing-docstring

import unittest

from taskforge.ql import Parser
from taskforge.ql.ast import AST, Expression
from taskforge.ql.tokens import Token

PARSER_TESTS = [{
    'name':
    'simple parse',
    'input':
    'milk and cookies',
    'output':
    AST(
        Expression(
            Token('and'),
            left=Expression(Token('milk')),
            right=Expression(Token('cookies'))), ),
}, {
    'name':
    'boolean parse',
    'input':
    'completed = false',
    'output':
    AST(
        Expression(
            Token('='),
            left=Expression(Token('completed')),
            right=Expression(Token('false'))), ),
}, {
    'name': 'simple all string parse',
    'input': 'milk -and cookies',
    'output': AST(Expression(Token('milk and cookies'))),
}, {
    'name':
    'complex parse',
    'input':
    '(priority > 5 and title ^ \'take out the trash\') or '
    '(context = "work" and (priority >= 2 or ("my little pony")))',
    'output':
    AST(
        Expression(
            Token('or'),
            right=Expression(
                Token('and'),
                left=Expression(
                    Token('='),
                    left=Expression(Token('context')),
                    right=Expression(Token('work'))),
                right=Expression(
                    Token('or'),
                    left=Expression(
                        Token('>='),
                        left=Expression(Token('priority')),
                        right=Expression(Token('2'))),
                    right=Expression(Token('my little pony'))),
            ),
            left=Expression(
                Token('and'),
                right=Expression(
                    Token('~'),
                    right=Expression(Token('take out the trash')),
                    left=Expression(Token('title'))),
                left=Expression(
                    Token('>'),
                    left=Expression(Token('priority')),
                    right=Expression(Token('5'))),
            ),
        ), ),
}, {
    'name':
    'parse bool infix',
    'input':
    'completed = false',
    'output':
    AST(
        Expression(
            Token('='),
            left=Expression(Token('completed')),
            right=Expression(Token('false'))))
}]


class ParserTests(unittest.TestCase):
    def test_parser(self):
        for test in PARSER_TESTS:
            with self.subTest(**test):
                parser = Parser(test['input'])
                ast = parser.parse()
                self.assertEqual(ast, test['output'])
