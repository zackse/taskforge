# pylint: disable=missing-docstring

import unittest
import os
from uuid import uuid1

from taskforge.lists.mongo import MongoDBList

from .list_utils import ListTests


class MongoDBListTests(unittest.TestCase, ListTests):
    def setUp(self):
        if not os.getenv('TASKFORGE_MONGO_TEST'):
            self.skipTest('TASKFORGE_MONGO_TEST not set')
        self.list = MongoDBList(db=uuid1().hex)

    def teardown(self):
        self.list._client.close()  # pylint: disable=protected-access
