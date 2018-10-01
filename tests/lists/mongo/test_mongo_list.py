# pylint: disable=missing-docstring

import unittest
from uuid import uuid1

import pytest

from task_forge.lists.mongo import List

from ..list_utils import ListTests, ListBenchmarks


@pytest.mark.slow
class MongoDBListTests(unittest.TestCase, ListTests):
    def setUp(self):
        self.list = List(db=uuid1().hex)

    def teardown(self):
        self.list._client.close()  # pylint: disable=protected-access


@pytest.mark.benchmark(group='MongoDB')
class TestMongoDBListPerformance(ListBenchmarks):
    @pytest.fixture
    def task_list(self):
        mongo = List(db=uuid1().hex)
        yield mongo
        mongo._client.close()  # pylint: disable=protected-access
