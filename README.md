# taskforge, a task management cli

Task management tool that integrates well with other services.

## Installation

### Install from Release

For now the only way to install taskforge is via pip:

```text
$ pip install taskforge-cli
```

## Usage


### task 

```text
usage: task [--version] <command> [<args>...]

A task management CLI that integrates with external services.

available commands:
   help                Print usage information about task commands
   add (new, a)        Add a new task to the list
   next (n)            Print the next or "current" task in the list
   todo                Print incomplete tasks in the list
   complete (done, d)  Complete tasks in the list.
   query (q, s, list)  Search or list tasks in the list

See 'task help <command>' for more information on a specific command.

```

### task add

```text
usage: task add [options] [<title>...]

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

```

### task complete

```text
usage: task complete [<ID>...]

Complete tasks by ID. If no IDs are provided then the current task indicated by
'task next' is completed.

```

### task next

```text
usage: task next [options]

Print the "next" or "current" task. This is calculated by the list as the
highest priority, oldest task in the list.

Default output format is:

$TASK_ID: $TASK_TITLE

You can modify the output with the options below.

options:
    -i, --id-only     Print only the task ID
    -t, --title-only  Print only the task title

```

### task query

```text
usage: task query [options] [<query>...]

Search or list tasks in this list.

QUERY will be concatenated using spaces and will be interpreted using the
Taskforge Query Language.

If no query is provided all tasks will be returned.

You can view information about the taskforge query
language using 'task help ql' or by visiting:

http://taskforge.io/docs/query_language

options:
    -o <format>, --output <format>  How to display the tasks which match the
                                    query. Available formats are: json, csv,
                                    table, text which are described below.
                                    [default: table]

output formats:

   Text output format is the same as for task next where each task will be
   printed on one line with the format:

   $TASK_ID: $TASK_TITLE

   Table output format lists tasks in an ascii table and it looks like this:

   | ID  | Created Date  | Completed Date  | Priority  | Title  | Context  |
   | --- | ------------- | --------------- | --------- | ------ | -------- |
   | $ID | $CREATED_DATE | $COMPLETED_DATE | $PRIORITY | $TITLE | $CONTEXT |

   JSON output format will "pretty print" a JSON array of the tasks in the list
   to stdout.  They will be properly indented 4 spaces and should be fairly
   human readable.  This is useful for migrating from one list implementation
   for another as you can redirect this output to a file then import it with:
   'task add --from-file $YOUR_JSON_FILE'

   CSV will output all task metadata in a csv format. It will write to stdout
   so you can use shell redirection to put it into a csv file like so:

   'task list --output csv > my_tasks.csv'

   This is useful for importing tasks into spreadsheet programs like Excel.

```

### task todo

```text
usage: task todo [options]

A convenience command for listing tasks which are incomplete.

options:
    -o <format>, --output <format>  How to display the tasks which match the
                                    query. Available formats are: json, csv,
                                    table, text. See 'task list --help' for
                                    more information on how each format is
                                    displayed. [default: table]


```

## Contributing

Contributions are greatly appreciated. We have a process for making a
contribution via Github that you can read in the 
[CONTRIBUTING.md](https://github.com/chasinglogic/taskforge/blob/master/CONTRIBUTING.md)
document.

## License

This code is distributed under the GNU General Public License

```text
    Copyright (C) 2018 Mathew Robinson

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
```