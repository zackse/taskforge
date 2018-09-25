# pylint: disable=missing-docstring

import unittest

from task_forge.ql.tokens import Token, Type


class TestTokens(unittest.TestCase):  # pylint: disable=too-many-public-methods
    def test_token_or(self):
        self.assertEqual(Token('or').token_type, Type.OR)

    def test_token_upper_or(self):
        self.assertEqual(Token('OR').token_type, Type.OR)

    def test_token_and(self):
        self.assertEqual(Token('and').token_type, Type.AND)

    def test_token_upper_and(self):
        self.assertEqual(Token('AND').token_type, Type.AND)

    def test_token_false(self):
        self.assertEqual(Token('false').token_type, Type.BOOLEAN)

    def test_token_upper_false(self):
        self.assertEqual(Token('False').token_type, Type.BOOLEAN)

    def test_token_true(self):
        self.assertEqual(Token('true').token_type, Type.BOOLEAN)

    def test_token_upper_true(self):
        self.assertEqual(Token('True').token_type, Type.BOOLEAN)

    def test_token_gt(self):
        self.assertEqual(Token('>').token_type, Type.GT)

    def test_token_lt(self):
        self.assertEqual(Token('<').token_type, Type.LT)

    def test_token_gte(self):
        self.assertEqual(Token('>=').token_type, Type.GTE)

    def test_token_lte(self):
        self.assertEqual(Token('<=').token_type, Type.LTE)

    def test_token_eq(self):
        self.assertEqual(Token('=').token_type, Type.EQ)

    def test_token_ne_shell(self):
        self.assertEqual(Token('^=').token_type, Type.NE)

    def test_token_ne(self):
        self.assertEqual(Token('!=').token_type, Type.NE)

    def test_token_like_shell(self):
        self.assertEqual(Token('^').token_type, Type.LIKE)

    def test_token_like(self):
        self.assertEqual(Token('~').token_type, Type.LIKE)

    def test_token_nlike_shell(self):
        self.assertEqual(Token('^^').token_type, Type.NLIKE)

    def test_token_nlike(self):
        self.assertEqual(Token('!~').token_type, Type.NLIKE)

    def test_token_lparen(self):
        self.assertEqual(Token('(').token_type, Type.LPAREN)

    def test_token_rparen(self):
        self.assertEqual(Token(')').token_type, Type.RPAREN)

    def test_token_num(self):
        self.assertEqual(Token('100').token_type, Type.NUMBER)

    def test_token_float(self):
        self.assertEqual(Token('1.00').token_type, Type.NUMBER)

    def test_token_string(self):
        self.assertEqual(Token('hello world').token_type, Type.STRING)

    def test_token_date_upper_pm(self):
        self.assertEqual(Token('2018-01-01 10:00 PM').token_type, Type.DATE)

    def test_token_date_pm(self):
        self.assertEqual(Token('2018-01-01 10:00 pm').token_type, Type.DATE)

    def test_token_date_upper_am(self):
        self.assertEqual(Token('2018-01-01 10:00 AM').token_type, Type.DATE)

    def test_token_date_am(self):
        self.assertEqual(Token('2018-01-01 10:00 am').token_type, Type.DATE)

    def test_token_date_24hr(self):
        self.assertEqual(Token('2018-01-01 10:00').token_type, Type.DATE)

    def test_token_date(self):
        self.assertEqual(Token('2018-01-01').token_type, Type.DATE)
