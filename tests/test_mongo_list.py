# pylint: disable=missing-docstring

import unittest
from uuid import uuid1

from taskforge.lists.mongo import MongoDBList

from .list_utils import ListTests


class SQLiteListTests(unittest.TestCase, ListTests):
    def setUp(self):
        self.list = MongoDBList(db=uuid1().hex)

    def teardown(self):
        self.list._client.close()  # pylint: disable=protected-access
