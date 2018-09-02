"""Implements the complete subcommand."""

import sys

from ..lists import NotFoundError
from ..task import Task
from .utils import inject_list


@inject_list
def complete_task(args, task_list=None):
    """Print the current task in task_list."""
    tasks = []
    if args.id:
        tasks = args.id
    else:
        try:
            current = task_list.current()
            tasks = [current.id]
        except NotFoundError:
            print('No ID given and no uncompleted task found')
            sys.exit(0)

    for task in tasks:
        task_list.complete(task)


def complete_cmd(parser):
    """Add the next command to parser."""
    sub_parser = parser.add_parser(
        'complete',
        aliases=['done', 'd'],
        help='Complete tasks in the list. If no ID given will'
        ' complete the current task.',
    )
    sub_parser.add_argument('id', metavar='ID', nargs='+', type=str)
    sub_parser.set_defaults(func=complete_task)
