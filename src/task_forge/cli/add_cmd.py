"""usage: task add [options] [<title>...]

Add or import tasks into the list.

options:
   -p <priority>, --priority <priority>  Create the task with the indicated
                                         priority, this can be an integer or
                                         float [default: 1.0]
   -b <body>, --body <body>              The body or "description" of the task
   -c <context>, --context <context>     The context in which to create the task

import options:
   -f <file>, --from-file <file>  Import tasks from the indicated JSON file

If an import option is provided all other options are ignored.
"""

import sys

from .utils import inject_list
from ..task import Task


@inject_list
def add_task(title, body='', context='default', priority=1.0, task_list=None):
    """Add a task to the configured task list"""
    task_list.add(Task(title, body=body, context=context, priority=priority))


@inject_list
def import_file(filename, task_list=None):
    """Import tasks from filename into the configured task list"""
    import json

    with open(filename) as tasks_file:
        task = json.load(tasks_file)
        if isinstance(task, list):
            tasks = [Task.from_dict(t) for t in task]
            task_list.add_multiple(tasks)
        else:
            task_list.add(task)


def run(args):
    """Parse the docopt args and call add_task."""
    if args['--from-file']:
        import_file(args['--from-file'])
        return

    if not args['<title>']:
        print('when not importing tasks title is required')
        sys.exit(1)

    title = ' '.join(args['<title>'])
    priority = float(args['--priority']) if args['--priority'] else 1.0
    context = args['--context'] if args['--context'] else 'default'
    body = args['--body'] if args['--body'] else ''
    add_task(title, body=body, context=context, priority=priority)
