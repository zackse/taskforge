"""Implements the todo subcommand."""

from ..ql.ast import AST, Expression
from ..ql.tokens import Token
from .query_cmd import print_tasks
from .utils import inject_list


@inject_list
def todo_task(args, task_list=None):
    """Print the current task in task_list."""
    ast = AST(
        Expression(
            Token('='),
            left=Expression(Token('completed')),
            right=Expression(Token('false'))))

    tasks = task_list.search(ast)
    print_tasks(tasks, output=args.output)


def todo_cmd(parser):
    """Add the next command to parser."""
    sub_parser = parser.add_parser(
        'todo',
        aliases=[],
        help='Print incomplete tasks in the list',
    )
    sub_parser.add_argument(
        '--output',
        '-o',
        type=str,
        default='table',
        choices=['text', 'table', 'json', 'csv'])
    sub_parser.add_argument('todo', metavar='TODO', nargs='*', type=str)
    sub_parser.set_defaults(func=todo_task)
