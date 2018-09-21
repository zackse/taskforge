# pylint: disable=missing-docstring

import pytest

from taskforge.ql.lexer import Lexer
from taskforge.ql.tokens import Token, Type


@pytest.mark.parametrize("query,expected", [(
    "milk and cookies",
    [
        Token('milk'),
        Token('and'),
        Token('cookies'),
    ],
), (
    "completed = false",
    [
        Token('completed'),
        Token('='),
        Token('false'),
    ],
), (
    "(priority > 0)",
    [
        Token('('),
        Token('priority'),
        Token('>'),
        Token('0'),
        Token(')'),
    ],
), (
    "milk -and cookies",
    [
        Token('milk'),
        Token('and', token_type=Type.STRING),
        Token('cookies'),
    ],
), (
    "(priority > 5 and title ^ \"take out the trash\") or "
    "(context = \"work\" and (priority >= 2 or (\"my little pony\")))",
    [
        Token('('),
        Token('priority'),
        Token('>'),
        Token('5'),
        Token('and'),
        Token('title'),
        Token('~'),
        Token('take out the trash'),
        Token(')'),
        Token('or'),
        Token('('),
        Token('context'),
        Token('='),
        Token('work'),
        Token('and'),
        Token('('),
        Token('priority'),
        Token('>='),
        Token('2'),
        Token('or'),
        Token('('),
        Token('my little pony'),
        Token(')'),
        Token(')'),
        Token(')'),
    ],
)])
def test_lexer(query, expected):
    lex = Lexer(query)
    tokens = list(lex)
    assert tokens == expected
    assert len(tokens) == len(expected)


@pytest.mark.slow
@pytest.mark.parametrize("query", [
    ('milk and cookies', ),
    ('milk -and cookies', ),
    ('completed = false', ),
    ('(priority > 5 and title ^ \'take out the trash\') or '
     '(context = "work" and (priority >= 2 or ("my little pony")))', ),
])
def test_lexer_performance(query, benchmark):
    """Benchmark the performance of various queries."""

    @benchmark
    def parse_query():  # pylint: disable=unused-variable
        """Benchmark query parsing"""
        lexer = Lexer(query)
        list(lexer)
