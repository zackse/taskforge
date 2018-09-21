# Pylint throws a ton of warnings about missing member variables that
# will be there when this class is subclassed
# pylint: skip-file

import pytest
import cProfile

from datetime import datetime

from taskforge.ql.ast import AST, Expression
from taskforge.ql.tokens import Token
from taskforge.ql.parser import Parser
from taskforge.task import Note, Task


class ListTests:
    """A class which implements the standard list tests"""

    @pytest.fixture
    def task_list(self):
        raise NotImplemented

    def test_add_one_and_find_by_id(self, task_list):
        task = Task("task 1")
        task_list.add(task)
        res = task_list.find_by_id(task.id)
        assert task == res

    def test_add_multiple(self, task_list):
        fixture = [
            Task("task 1"),
            Task("task 2"),
            Task("task 3"),
        ]

        task_list.add_multiple(fixture)
        result = task_list.list()
        assert fixture == result

    def test_complete_a_task(self, task_list):
        task = Task('task to complete')
        task_list.add(task)
        task_list.complete(task.id)
        result = task_list.find_by_id(task.id)
        assert result.is_completed()

    def test_return_correct_current_task(self, task_list):
        tasks = [
            Task("task 1"),
            Task("task 2"),
        ]

        task_list.add_multiple(tasks)
        current = task_list.current()
        assert tasks[0] == current
        task_list.complete(tasks[0].id)
        current = task_list.current()
        assert tasks[1] == current

    def test_add_a_note_to_a_task(self, task_list):
        task = Task("task to be noted")
        task_list.add(task)
        note = Note("a note")
        task_list.add_note(task.id, note)
        noted = task_list.find_by_id(task.id)
        assert noted.notes == [note]

    def test_update_a_task(self, task_list):
        task = Task("task to update")
        task_list.add(task)
        to_update = task_list.find_by_id(task.id)
        to_update.title = "task updated"
        task_list.update(to_update)
        updated = task_list.find_by_id(task.id)
        assert updated.title == "task updated"

    def run_query_test(self,
                       task_list=None,
                       query='',
                       fixture=None,
                       expected=None):
        ast = Parser(query).parse()
        task_list.add_multiple(fixture)
        result = task_list.search(ast)
        assert result == expected

    def test_query_simple_title(self, task_list):
        fixture = [
            Task("task 1"),
            Task("task 2"),
        ]
        self.run_query_test(
            task_list=task_list,
            query="title = \"task 1\"",
            fixture=fixture,
            expected=[fixture[0]],
        )

    def test_query_other_context(self, task_list):
        fixture = [
            Task("task 1"),
            Task("task 2"),
            Task("other task", context="other"),
            Task("task 3"),
        ]
        self.run_query_test(
            task_list=task_list,
            query="context = other",
            fixture=fixture,
            expected=[fixture[2]],
        )

    def test_query_multiple_ors(self, task_list):
        fixture = [
            Task("task 1"),
            Task("task 2"),
            Task("other task", context="other"),
            Task("task 3"),
            Task("task 4"),
        ]
        self.run_query_test(
            task_list=task_list,
            query=
            "title = \"task 4\" or title = \"task 1\" or title = \"other task\"",
            fixture=fixture,
            expected=[fixture[0], fixture[4], fixture[2]],
        )

    def test_query_grouped_expressions(self, task_list):
        fixture = [
            Task("task 1"),
            Task("task 2"),
            Task("other task", context="other"),
            Task("task 3"),
            Task("task 4"),
        ]
        self.run_query_test(
            task_list=task_list,
            query=
            "(title = \"task 1\" and context = \"default\") or (context = \"other\")",
            fixture=fixture,
            expected=[fixture[0], fixture[2]],
        )

    def test_query_string_literal_only(self, task_list):
        fixture = [
            Task("task 1"),
            Task("task 2"),
            Task("other task", context="other"),
            Task("task 3"),
            Task("task 4"),
        ]
        self.run_query_test(
            task_list=task_list,
            query="task",
            fixture=fixture,
            expected=[
                fixture[0], fixture[1], fixture[2], fixture[3], fixture[4]
            ])

    def test_query_priority_equals_1_0(self, task_list):
        fixture = [
            Task("task 1"),
            Task("task 2"),
            Task("other task", context="other"),
            Task("task 3"),
            Task("task 4"),
        ]
        self.run_query_test(
            task_list=task_list,
            query="priority = 1.0",
            fixture=fixture,
            expected=[
                fixture[0], fixture[1], fixture[2], fixture[3], fixture[4]
            ])

    def test_query_priority_gt_1(self, task_list):
        fixture = [
            Task("task 1"),
            Task("task 2"),
            Task("other task", context="other", priority=2.0),
            Task("task 3"),
            Task("task 4"),
        ]
        self.run_query_test(
            task_list=task_list,
            query="priority > 1",
            fixture=fixture,
            expected=[fixture[2]])

    def test_completed_false(self, task_list):
        fixture = [
            Task("task 1"),
            Task("task 2", completed_date=datetime.now()),
            Task("other task", context="other", priority=2.0),
            Task("task 3", completed_date=datetime.now()),
            Task("task 4"),
        ]

        self.run_query_test(
            task_list=task_list,
            query="completed = false",
            fixture=fixture,
            expected=[fixture[0], fixture[2], fixture[4]])

    def test_completed_true(self, task_list):
        fixture = [
            Task("task 1"),
            Task("task 2", completed_date=datetime.now()),
            Task("other task", context="other", priority=2.0),
            Task("task 3", completed_date=datetime.now()),
            Task("task 4"),
        ]

        self.run_query_test(
            task_list=task_list,
            query="completed = true",
            fixture=fixture,
            expected=[fixture[1], fixture[3]])


#     @pytest.fixture(autouse=True)
#     def test_query_benchmark(self, benchmark):
#         # Hand-crafted artisinal Abstract Syntax Tree
#         ast = AST(
#             Expression(
#                 Token('or'),
#                 right=Expression(
#                     Token('and'),
#                     left=Expression(
#                         Token('='),
#                         left=Expression(Token('context')),
#                         right=Expression(Token('work'))),
#                     right=Expression(
#                         Token('or'),
#                         left=Expression(
#                             Token('>='),
#                             left=Expression(Token('priority')),
#                             right=Expression(Token('2'))),
#                         right=Expression(Token('my little pony'))),
#                 ),
#                 left=Expression(
#                     Token('and'),
#                     right=Expression(
#                         Token('~'),
#                         right=Expression(Token('take out the trash')),
#                         left=Expression(Token('title'))),
#                     left=Expression(
#                         Token('>'),
#                         left=Expression(Token('priority')),
#                         right=Expression(Token('5'))),
#                 ),
#             ), )

#         tasks = [
#             Task("my little pony"),
#             Task("this task won't match anything"),
#             Task("a priority 2 task", priority=2.0),
#             Task("take out the trash", priority=5.0),
#             Task("work task 1", context="work"),
#             Task("work task 2", context="work"),
#             Task("task 1"),
#             Task("task 2"),
#             Task("task 3"),
#             Task("task 4"),
#         ]

#         task_list.add_multiple(tasks)
#         benchmark(task_list.search, args=(ast, ))
