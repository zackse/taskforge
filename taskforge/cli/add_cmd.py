"""Implements the new subcommand."""

import json

from ..task import Task
from .utils import inject_list


@inject_list
def add_task(args, task_list=None):
    """Print the current task in task_list."""
    if args.from_file:
        with open(args.from_file) as tasks_file:
            task = json.load(tasks_file)
            if isinstance(task, list):
                tasks = [Task.from_dict(t) for t in task]
                task_list.add_multiple(tasks)
            else:
                task_list.add(task)
            return

    task_list.add(
        Task(
            ' '.join(args.title),
            body=args.body,
            context=args.context,
            priority=args.priority))


def add_cmd(parser):
    """Add the next command to parser."""
    sub_parser = parser.add_parser(
        'add',
        aliases=['new', 'a'],
        help='Add a new task to the list',
    )
    sub_parser.add_argument(
        '--from-file',
        '-f',
        type=str,
        help='A JSON file which to load tasks from, '
        'if provided all other arguments are ignored.')
    sub_parser.add_argument('--priority', '-p', type=float, default=1.0)
    sub_parser.add_argument('--context', '-c', type=str, default='default')
    sub_parser.add_argument('--body', '-b', type=str)
    sub_parser.add_argument('title', metavar='TITLE', nargs='*', type=str)
    sub_parser.set_defaults(func=add_task)
