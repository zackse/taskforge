# pylint: disable=missing-docstring

import unittest
from tempfile import NamedTemporaryFile

from taskforge.lists.sqlite import SQLiteList
from taskforge.task import Task

from .list_utils import ListTests


class SQLiteListTests(unittest.TestCase, ListTests):
    def setUp(self):
        self.tmpfile = NamedTemporaryFile()
        self.list = SQLiteList(file_name=self.tmpfile.name, create_tables=True)

    def teardown(self):
        self.tmpfile.close()

    def test_save_and_load(self):
        tasks = [Task('test 1'), Task('test 2')]
        self.list.add_multiple(tasks)
        new_list = SQLiteList(file_name=self.tmpfile.name)
        self.assertCountEqual(new_list.list(), tasks)
