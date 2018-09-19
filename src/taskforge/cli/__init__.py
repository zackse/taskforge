"""Entry point for the Taskforge CLI."""

import argparse

from .add_cmd import add_cmd
from .complete_cmd import complete_cmd
# from .edit_cmd import edit_cmd
from .next_cmd import next_cmd
from .query_cmd import query_cmd
from .todo_cmd import todo_cmd


def main():
    """Entry point function for the Taskforge CLI."""
    parser = argparse.ArgumentParser(prog='task')
    subparsers = parser.add_subparsers()
    next_cmd(subparsers)
    add_cmd(subparsers)
    query_cmd(subparsers)
    todo_cmd(subparsers)
    complete_cmd(subparsers)
    # edit_cmd(subparsers)

    args = parser.parse_args()
    if hasattr(args, 'func'):
        args.func(args)
    else:
        parser.print_help()


if __name__ == '__main__':
    main()
