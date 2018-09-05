# pylint: disable=missing-docstring

import unittest
from datetime import datetime

from taskforge.ql.ast import Expression
from taskforge.ql.tokens import Token


class ExpressionTests(unittest.TestCase):
    def test_expression_values_literals(self):
        literals = [{
            'literal': '1.0',
            'value': 1.0,
        }, {
            'literal': 'hello world',
            'value': 'hello world',
        }, {
            'literal': '2018-01-01',
            'value': datetime(year=2018, month=1, day=1)
        }, {
            'literal': 'True',
            'value': True,
        }, {
            'literal': 'true',
            'value': True,
        }, {
            'literal': 'False',
            'value': False,
        }, {
            'literal': 'false',
            'value': False,
        }]

        for literal in literals:
            with self.subTest(**literal):
                exp = Expression(Token(literal['literal']))
                self.assertEqual(type(exp.value), type(literal['value']))
                self.assertEqual(exp.value, literal['value'])

    def test_is_infix(self):
        infix = Expression(
            Token('='), left=Token('title'), right=Token('right'))
        self.assertTrue(infix.is_infix())

    def test_is_literal(self):
        literal = Expression(Token('milk'))
        self.assertTrue(literal.is_literal())

    def test_date_formats(self):
        date_strings = [
            {
                'date_string': '2018-01-01',
                'expected': datetime(year=2018, month=1, day=1)
            },
            {
                'date_string': '2018-01-01 01:01',
                'expected': datetime(
                    year=2018, month=1, day=1, hour=1, minute=1)
            },
            {
                'date_string':
                '2018-01-01 01:01:10',
                'expected':
                datetime(
                    year=2018, month=1, day=1, hour=1, minute=1, second=10)
            },
            {
                'date_string':
                '2018-01-01 01:01:10 PM',
                'expected':
                datetime(
                    year=2018, month=1, day=1, hour=13, minute=1, second=10)
            },
            {
                'date_string':
                '2018-01-01 01:01:10PM',
                'expected':
                datetime(
                    year=2018, month=1, day=1, hour=13, minute=1, second=10)
            },
            {
                'date_string': '2018-01-01 01:01 PM',
                'expected': datetime(
                    year=2018, month=1, day=1, hour=13, minute=1)
            },
            {
                'date_string': '2018-01-01 01:01PM',
                'expected': datetime(
                    year=2018, month=1, day=1, hour=13, minute=1)
            },
        ]

        for date in date_strings:
            with self.subTest(**date):
                result = Expression.parse_date(date['date_string'])
                self.assertEqual(result, date['expected'])
