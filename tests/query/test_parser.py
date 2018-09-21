# pylint: disable=missing-docstring

import pytest

from taskforge.ql.ast import AST, Expression
from taskforge.ql.parser import Parser
from taskforge.ql.tokens import Token


@pytest.mark.parametrize("query,ast", [(
    'milk and cookies',
    AST(
        Expression(
            Token('and'),
            left=Expression(Token('milk')),
            right=Expression(Token('cookies'))), ),
), (
    'completed = false',
    AST(
        Expression(
            Token('='),
            left=Expression(Token('completed')),
            right=Expression(Token('false'))), ),
), (
    'milk -and cookies',
    AST(Expression(Token('milk and cookies'))),
), (
    '(priority > 5 and title ^ \'take out the trash\') or '
    '(context = "work" and (priority >= 2 or ("my little pony")))',
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
), ('completed = false',
    AST(
        Expression(
            Token('='),
            left=Expression(Token('completed')),
            right=Expression(Token('false')))))])
def test_parser(query, ast):
    parser = Parser(query)
    assert parser.parse() == ast


@pytest.mark.slow
@pytest.mark.parametrize("query", [
    ('milk and cookies', ),
    ('milk -and cookies', ),
    ('completed = false', ),
    ('(priority > 5 and title ^ \'take out the trash\') or '
     '(context = "work" and (priority >= 2 or ("my little pony")))', ),
])
def test_parser_performance(query, benchmark):
    """Benchmark the performance of various queries."""

    @benchmark
    def parse_query():  # pylint: disable=unused-variable
        """Benchmark query parsing"""
        parser = Parser(query)
        parser.parse()
