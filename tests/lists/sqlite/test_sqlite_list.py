# pylint: disable=missing-docstring

import pytest

from taskforge.lists.sqlite import List
from taskforge.task import Task

from ..list_utils import ListTests


class TestSQLiteList(ListTests):
    @pytest.fixture
    def task_list(self, tmpdir):
        tmpfile = tmpdir.join("tasks.sqlite3")
        return List(file_name=tmpfile, create_tables=True)
