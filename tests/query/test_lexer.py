# pylint: disable=missing-docstring

import unittest

from taskforge.ql.lexer import Lexer
from taskforge.ql.tokens import Token, Type

LEXER_TESTS = [{
    'name': "simple lex",
    'input': "milk and cookies",
    'expected': [
        Token('milk'),
        Token('and'),
        Token('cookies'),
    ],
}, {
    'name': "boolean lex",
    'input': "completed = false",
    'expected': [
        Token('completed'),
        Token('='),
        Token('false'),
    ],
}, {
    'name':
    "single grouped expression",
    'input':
    "(priority > 0)",
    'expected': [
        Token('('),
        Token('priority'),
        Token('>'),
        Token('0'),
        Token(')'),
    ],
}, {
    'name':
    "keyword excaped lex",
    'input':
    "milk -and cookies",
    'expected': [
        Token('milk'),
        Token('and', token_type=Type.STRING),
        Token('cookies'),
    ],
}, {
    'name':
    "complicated lex",
    'input':
    "(priority > 5 and title ^ \"take out the trash\") or "
    "(context = \"work\" and (priority >= 2 or (\"my little pony\")))",
    'expected': [
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
}]


class LexerTests(unittest.TestCase):
    def test_lexer_tokens(self):
        """Test Lexer"""
        for test in LEXER_TESTS:
            with self.subTest(name=test['name'], query=test['input']):
                self.run_lexer_test(test)

    def run_lexer_test(self, test):
        lex = Lexer(test['input'])
        tokens = list(lex)
        expected = test['expected']
        self.assertEqual(len(tokens), len(expected))
        for token, expected in zip(tokens, expected):
            self.assertEqual(token.__dict__, expected.__dict__)
