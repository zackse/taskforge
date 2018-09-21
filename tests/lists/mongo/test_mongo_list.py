# pylint: disable=missing-docstring

import unittest
import cProfile
import os
from uuid import uuid1

from taskforge.lists.mongo import List
from taskforge.ql.ast import AST, Expression
from taskforge.ql.tokens import Token

from ..list_utils import ListTests


class MongoDBListTests(unittest.TestCase, ListTests):
    def setUp(self):
        if not os.getenv('TASKFORGE_MONGO_TEST'):
            self.skipTest('TASKFORGE_MONGO_TEST not set')
        self.list = List(db=uuid1().hex)

    def teardown(self):
        self.list._client.close()  # pylint: disable=protected-access
