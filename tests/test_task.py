import unittest

from taskforge.task import Task


class TaskTests(unittest.TestCase):

    def test_unique_ids(self):
        t1 = Task('task 1')
        t2 = Task('task 2')
        t3 = Task('task 3')
        self.assertNotEqual(t1, t2)
        self.assertNotEqual(t1, t3)
        self.assertNotEqual(t2, t3)
        self.assertNotEqual(t1.created_date, t3.created_date)

    def test_sort_order(self):
        t1 = Task('task 1')
        t2 = Task('task 2')
        t3 = Task('task 3')

        list1 = sorted([
            t3,
            t2,
            t1
        ])

        self.assertEqual(list1[0], t1)
        self.assertEqual(list1[1], t2)
        self.assertEqual(list1[2], t3)

        t1.priority = 3.0
        t2.priority = 1.0
        t3.priority = 2.0

        list2 = sorted([
            t3,
            t2,
            t1
        ])

        self.assertEqual(list2[0], t1)
        self.assertEqual(list2[1], t3)
        self.assertEqual(list2[2], t2)
