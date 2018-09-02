# pylint: disable=missing-docstring

import unittest

from taskforge.ql.tokens import Token, Type


class TestTokens(unittest.TestCase):
    def test_token_types(self):
        self.assertEqual(Token('or').token_type, Type.OR)
        self.assertEqual(Token('OR').token_type, Type.OR)
        self.assertEqual(Token('and').token_type, Type.AND)
        self.assertEqual(Token('AND').token_type, Type.AND)
        self.assertEqual(Token('false').token_type, Type.BOOLEAN)
        self.assertEqual(Token('False').token_type, Type.BOOLEAN)
        self.assertEqual(Token('true').token_type, Type.BOOLEAN)
        self.assertEqual(Token('True').token_type, Type.BOOLEAN)
        self.assertEqual(Token('>').token_type, Type.GT)
        self.assertEqual(Token('<').token_type, Type.LT)
        self.assertEqual(Token('>=').token_type, Type.GTE)
        self.assertEqual(Token('<=').token_type, Type.LTE)
        self.assertEqual(Token('=').token_type, Type.EQ)
        self.assertEqual(Token('^=').token_type, Type.NE)
        self.assertEqual(Token('!=').token_type, Type.NE)
        self.assertEqual(Token('^').token_type, Type.LIKE)
        self.assertEqual(Token('~').token_type, Type.LIKE)
        self.assertEqual(Token('^^').token_type, Type.NLIKE)
        self.assertEqual(Token('!~').token_type, Type.NLIKE)
        self.assertEqual(Token('(').token_type, Type.LPAREN)
        self.assertEqual(Token(')').token_type, Type.RPAREN)
        self.assertEqual(Token('100').token_type, Type.NUMBER)
        self.assertEqual(Token('hello world').token_type, Type.STRING)
        self.assertEqual(Token('2018-01-01 10:00 PM').token_type, Type.DATE)
        self.assertEqual(Token('2018-01-01 10:00 pm').token_type, Type.DATE)
        self.assertEqual(Token('2018-01-01 10:00 AM').token_type, Type.DATE)
        self.assertEqual(Token('2018-01-01 10:00 am').token_type, Type.DATE)
        self.assertEqual(Token('2018-01-01 10:00').token_type, Type.DATE)
        self.assertEqual(Token('2018-01-01').token_type, Type.DATE)
