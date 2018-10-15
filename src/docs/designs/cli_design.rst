task (CLI Client)
=================

Goals
-----

-  Task management CLI as the first usable frontend

   -  Supports use of different lists via a configuration file
   -  Follows Unix best practices
   -  Uses a subcommand interface the following commands will be
      supported:

      -  `new <#new-subcommand>`_
      -  `note <#note-subcommand>`_
      -  `complete <#complete-subcommand>`_
      -  `query <#query-subcommand>`_
      -  `next <#next-subcommand>`_
      -  `edit <#edit-subcommand>`_

Design
------

Configuration File
~~~~~~~~~~~~~~~~~~

The task CLI will be configured with TOML.

An example config file will look like:

.. code::

   [list]
   name = "sqlite"

   [list.config]
   file = "~/.taskforge.d/tasks.sqlite3"

To start there will be two sections, ``[list]`` which has a single key
name. This name corresponds to the list implementation the user wants to load.
``[list.confg]`` is a section filled with arbitrary key value pairs that are
passed to the constructor of the list implementation as kwargs deconstructed
using the ``**`` operator.

New Subcommand
--------------

The new subcommand will accept the following flags:

-  ``--body TEXT_BODY`` populates the task body
-  ``--priority PRIORITY_NUMBER`` populates the task priority
-  ``--context CONTEXT`` populates the task context
-  ``--from-file PATH_TO_FILE`` loads a task/s from a yaml file or csv
   file

It takes VarArgs and concatenates them into the title of a new task. So that:

.. code:: bash

   task new write a design doc

Will create a task with the title “write a design doc." Flags as
described above can be passed to populate other fields of the Task.
Otherwise the flag fields will get the defaults described below:

-  body: None
-  priority: 0.0
-  context: “default”

The VarArgs are ignored if ``--from-file`` is provided. If the file is a
.csv then new will assume it is a CSV with the following format:

.. code:: text

   title,body,context,priority
   record_title,record_body,record_context,record_priority

The order of the columns is not important. Only the title, priority, and
context columns are required. Values can be omitted for the optional
comments for any record which does require them.

Note Subcommand
---------------

The note subcommand takes no flags and one argument: the ID of the task
to add a note to. So that:

.. code:: bash

   task note TASK_ID

Opens up your ``$EDITOR`` and allows you to input text that will then be
used as the body of a note which is attached to the task.

Complete Subcommand
-------------------

The complete subcommand takes no flags and one argument: the ID of the
task to add a note to. So that:

.. code:: bash

   task complete TASK_ID

Will complete the task indicated by ``TASK_ID``.

Query Subcommand
----------------

The query subcommand takes the following flags:

-  ``--completed`` a convenience flag to show completed results
-  ``--csv`` print results as a CSV
-  ``--raw`` print no decoration on task table (i.e. remove the “\|” and
   “-” characters)
-  ``--id-only`` print only matching task ID’s

It takes VarArgs and concatenates them into a query using the `Query
Language <#query-language>`__ parser. It then prints each task in a
table using the following format:

.. code:: text

   --------------------------------------------
   | ID      | Created Date      | Title      |
   --------------------------------------------
   | TASK_ID | TASK_CREATED_DATE | TASK_TITLE |
   --------------------------------------------

If raw is given:

.. code:: text

   ID      Created Date      Title
   TASK_ID TASK_CREATED_DATE TASK_TITLE

If ID is given only a newline separated list of
TASK_IDs are printed with no headers.

Next Subcommand
---------------

The next subcommand takes the following flags:

-  ``--title-only`` print only the task title
-  ``--id-only`` print only the task id

But it takes no arguments. It returns the item currently at the “top” of
the list (sorted by oldest date and highest priority). It prints it like
so:

.. code:: text

   TASK_ID TASK_CREATED_DATE TASK_TITLE

If title or id only flags are given then only that field is printed.

Edit Subcommand
---------------

The edit subcommand takes one argument: the task ID. It opens the
indicated task in ``$EDITOR`` as a yaml file and includes all fields
from the task. Upon saving and exiting the file will be read, parsed,
and the task will be updated with that info.

Future Work / Ideas
-------------------

-  Configurable canned queries
