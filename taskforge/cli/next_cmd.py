"""Implements the next subcommand."""

from ..lists import NotFoundError
from .utils import inject_list


@inject_list
def print_next(args, task_list=None):
    """Print the current task in task_list."""
    try:
        task = task_list.current()
    except NotFoundError:
        print('No current task!')
        return

    if args.title_only:
        print(task.title)
    elif args.id_only:
        print(task.id)
    else:
        print(task)


def next_cmd(parser):
    """Add the next command to parser."""
    sub_parser = parser.add_parser(
        'next', aliases=['n'], help='Print the next task in the list')
    sub_parser.add_argument('--title-only', '-t', action='store_true')
    sub_parser.add_argument('--id-only', '-i', action='store_true')
    sub_parser.set_defaults(func=print_next)
