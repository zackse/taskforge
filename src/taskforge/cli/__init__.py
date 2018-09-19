"""usage: task [--version] <command> [<args>...]

A task management CLI that integrates with external services.

available commands:
   help                Print usage information about task commands
   add (new, a)        Add a new task to the list
   next (n)            Print the next or "current" task in the list
   todo                Print incomplete tasks in the list
   complete (done, d)  Complete tasks in the list.
   query (q, s, list)  Search or list tasks in the list

See 'task help <command>' for more information on a specific command.
"""

import sys
from importlib import import_module

from docopt import docopt

ALIASES = {
    'n': 'next',
    'new': 'add',
    'a': 'add',
    'd': 'complete',
    'done': 'complete',
    'q': 'query',
    's': 'query',
    'list': 'query',
}


def main():
    """CLI entrypoint, handles subcommand parsing"""
    args = docopt(__doc__, version='task version 0.1.0', options_first=True)
    if not args['<command>']:
        print(__doc__)
        sys.exit(1)

    command = args['<command>']
    try:
        if command == 'help':
            if args['<args>']:
                command_mod = import_module('taskforge.cli.{}_cmd'.format(
                    args['<args>'][0]))
                print(command_mod.__doc__)
            else:
                print(__doc__)
            sys.exit(0)

        command = ALIASES.get(command, command)
        command_mod = import_module('taskforge.cli.{}_cmd'.format(command))
        argv = [command] + args['<args>']
        command_mod.run(docopt(command_mod.__doc__, argv=argv))
        sys.exit(0)
    except ImportError:
        print('{} is not a known task command.'.format(command))
        sys.exit(1)


if __name__ == '__main__':
    main()
