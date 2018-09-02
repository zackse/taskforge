"""
Implements the query subcommand
"""

import csv
import json
import sys

from ..ql import Parser
from ..task import DATE_FORMAT
from .utils import inject_list


def print_table(tasks):
    """Print an ASCII table of the tasks"""
    rows = [
        [
            'ID',
            'Created Date',
            'Completed Date',
            'Priority',
            'Title',
            'Context'
        ]
    ]
    rows += [
        [
            task.id,
            task.created_date,
            task.completed_date,
            task.priority,
            task.title,
            task.context
        ]
        for task in tasks
    ]

    wcolumns = None
    for columns in rows:
        if not wcolumns:
            wcolumns = [len(str(x)) for x in columns]
        else:
            wcolumns = [max(x, len(str(y))) for x, y in zip(wcolumns, columns)]

    # print columns with the maximum width
    for columns in rows:
        cols = [str(c).ljust(w) for w, c in zip(wcolumns, columns)]
        print("| {} |".format(" | ".join(list(cols))))


def print_text(tasks):
    """Print the __repr__ of all tasks in the list"""
    for task in tasks:
        print(task)


def print_json(tasks):
    """Print a list of tasks as json to stdout"""
    dicts = [task.to_dict() for task in tasks]
    json.dump(dicts, sys.stdout, indent="\t")
    # add a newline to output
    print()


def print_csv(tasks):
    """Print a list of tasks as csv to stdout"""
    writer = csv.DictWriter(
        sys.stdout,
        extrasaction='ignore',
        fieldnames=[
            'id',
            'created_date',
            'completed_date',
            'priority',
            'title',
            'context',
            'body',
        ]
    )

    writer.writeheader()
    for task in tasks:
        writer.writerow(task.to_dict())

def print_tasks(tasks, output='table'):
    """Print tasks using the print function which corresponds to output"""
    if output == 'table':
        print_table(tasks)
    elif output == 'text':
        print_text(tasks)
    elif output == 'json':
        print_json(tasks)
    elif output == 'csv':
        print_csv(tasks)
    else:
        print('{} is not a valid output format'.format(args.output))


@inject_list
def query_task(args, task_list=None):
    """Print the current task in task_list"""
    if args.query:
        ast = Parser(' '.join(args.query)).parse()
        tasks = task_list.search(ast)
    else:
        tasks = task_list.list()

    print_tasks(tasks, output=args.output)


def query_cmd(parser):
    """Add the next command to parser"""
    sub_parser = parser.add_parser(
        'query',
        aliases=['q', 's', 'search', 'list'],
        help='Query tasks in the list.',
    )
    sub_parser.add_argument(
        '--output', '-o',
        type=str,
        default='table',
        choices=['text', 'table', 'json', 'csv']
    )
    sub_parser.add_argument('query', metavar='QUERY', nargs='*', type=str)
    sub_parser.set_defaults(func=query_task)
