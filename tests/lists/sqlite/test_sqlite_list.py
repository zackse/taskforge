# pylint: disable=missing-docstring

import pytest
import unittest

from tempfile import NamedTemporaryFile

from task_forge.lists.sqlite import List
from task_forge.task import Task

from ..list_utils import ListTests, ListBenchmarks


class SQLiteListTests(unittest.TestCase, ListTests):
    def setUp(self):
        self.tmpfile = NamedTemporaryFile()
        self.list = List(file_name=self.tmpfile.name, create_tables=True)

    def tearDown(self):
        self.tmpfile.close()


@pytest.mark.benchmark(group='SQLite')
class TestSQLiteListPerformance(ListBenchmarks):
    @pytest.fixture
    def task_list(self, tmpdir):
        tmpfile = tmpdir.join("tasks.sqlite3")
        return List(file_name=tmpfile, create_tables=True)
