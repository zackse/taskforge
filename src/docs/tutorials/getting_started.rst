Getting Started
===============

Welcome to the Taskforge getting started guide. This is everything you need to
know to get Taskforge up and running.

.. contents::
   :local:

Prerequisites
-------------

Before you install Taskforge you'll need a few things installed. Writing an
installation guide for all the prerequisites is outside the scope of this
document so we will simply link to the canonical documentation below:

- `Python 3 <https://python.org>`_
- `PIP (Included with Python versions 3.4 and above)
  <https://pip.pypa.io/en/stable/installing/>`_

Certain list implementations will require additional packages not installed in
this guide to work. See the :doc:`../lists/index` documentation for your preferred list to
know more. In this guide we will use the SQLite list because it will work on
most platforms with no additional setup.

Installing
----------

First, let's install Taskforge. At the time of this writing the only method for
installing Taskforge is from pip:

.. console::

   $ pip3 install taskforge-cli

.. note:: The pip command can vary slightly based on your platform, ``pip3`` is
   used here because it will work on most platforms.

   For example Windows users will need to do:
   
   .. console::

      $ python.exe -m pip install taskforge-cli


Using Taskforge
---------------

Your first task
+++++++++++++++

Now that Taskforge is installed we can start using it. Out of the box taskforge
will use a SQLite database to store and retrieve tasks. Lets add a task now:

.. console::

   $ task add complete the Taskforge tutorial
   $


To see what tasks are in our list we can use ``task list``. Let's run it now:

.. console::

   $ task list
   | ID                               | Created Date               | Completed Date | Priority | Title                           | Context |
   | -------------------------------- | -------------------------- | -------------- | -------- | ------------------------------- | ------- |
   | eabdeee413ef442fa68c994119d817d2 | 2018-09-23 18:41:18.858741 | None           | 1.0      | complete the taskforge tutorial | default |
   $

What's next?
++++++++++++

If we want to see what our current task is you can use ``task next`` or 
``task current``:

.. console::

   $ task next
   eabdeee413ef442fa68c994119d817d2: complete the taskforge tutorial
   $

Taskforge defines the 'current' task as the highest priority task. If all tasks
are of equal priority then the 'current' task is the one with the oldest created
date. To demonstrate let's add a few more tasks: 

.. console::

   $ task add another default priority task
   $ task add --priority 2 a high priority task

This introduces a new flag ``--priority``. You can set many fields on a task via
flags to the add command. See the :doc:`../cli/task_add` documentation for more
information.

Now our ``task list`` should look like this:

.. console::

   $ task list
   | ID                               | Created Date               | Completed Date | Priority | Title                           | Context |
   | -------------------------------- | -------------------------- | -------------- | -------- | ------------------------------- | ------- |
   | eabdeee413ef442fa68c994119d817d2 | 2018-09-23 18:41:18.858741 | None           | 1.0      | complete the taskforge tutorial | default |
   | 1e634ced06d64093a747f38da024f9a6 | 2018-09-23 18:46:05.198426 | None           | 1.0      | another default priority task   | default |
   | 265b67ff298643dbb05950f3394a5ab0 | 2018-09-23 18:46:30.082289 | None           | 2.0      | a high priority task            | default |
   $

If we run ``task next`` now we'll see that the 'a high priority task' is the
current task:

.. console::

   $ task next
   265b67ff298643dbb05950f3394a5ab0: a high priority task
   $

This is because priority, in the Taskforge world, is the #1 indicator of what
you should be working on. Then you should be working on whatever has been
waiting the longest.

Completing tasks
++++++++++++++++

You can complete tasks with ``task done`` or ``task complete``. Let's complete
our high priority task:

.. console::
   
   $ task next
   265b67ff298643dbb05950f3394a5ab0: a high priority task
   $ task done 265b67ff298643dbb05950f3394a5ab0
   $

Every task has a unique ID. Most commands will show you this ID for easy with
other commands like done which take a Task ID as an argument. 

Viewing incomplete tasks
++++++++++++++++++++++++

Now that we've completed this task we'll see that the current task has changed:

.. console::

   $ task next
   eabdeee413ef442fa68c994119d817d2: complete the taskforge tutorial
   $

However if we run ``task list`` we will still see the completed task:

.. console::

   $ task list
   | ID                               | Created Date               | Completed Date             | Priority | Title                           | Context |
   | -------------------------------- | -------------------------- | -------------------------- | -------- | ------------------------------- | ------- |
   | eabdeee413ef442fa68c994119d817d2 | 2018-09-23 18:41:18.858741 | None                       | 1.0      | complete the taskforge tutorial | default |
   | 1e634ced06d64093a747f38da024f9a6 | 2018-09-23 18:46:05.198426 | None                       | 1.0      | another default priority task   | default |
   | 265b67ff298643dbb05950f3394a5ab0 | 2018-09-23 18:46:30.082289 | 2018-09-23 18:55:24.277754 | 2.0      | a high priority task            | default |
   $


As your task list grows finding tasks that need to be done using ``task list``
can be overwhelming. Luckily, Taskforge has a :doc:`../how_tos/query_language` we can use to
search tasks. See the linked documentation for full instructions, for our
purposes we simply need to run the following:

.. console::

   $ task query completed = false
   | ID                               | Created Date               | Completed Date | Priority | Title                           | Context |
   | -------------------------------- | -------------------------- | -------------- | -------- | ------------------------------- | ------- |
   | eabdeee413ef442fa68c994119d817d2 | 2018-09-23 18:41:18.858741 | None           | 1.0      | complete the taskforge tutorial | default |
   | 1e634ced06d64093a747f38da024f9a6 | 2018-09-23 18:46:05.198426 | None           | 1.0      | another default priority task   | default |
   $


This shows us all tasks which are incomplete. This is such a common query that
there is a shortcut command for displaying this information ``task todo``:

.. console::

   $ task todo
   | ID                               | Created Date               | Completed Date | Priority | Title                           | Context |
   | -------------------------------- | -------------------------- | -------------- | -------- | ------------------------------- | ------- |
   | eabdeee413ef442fa68c994119d817d2 | 2018-09-23 18:41:18.858741 | None           | 1.0      | complete the taskforge tutorial | default |
   | 1e634ced06d64093a747f38da024f9a6 | 2018-09-23 18:46:05.198426 | None           | 1.0      | another default priority task   | default |
   $


Re-ordering tasks
+++++++++++++++++

Sometimes a task which you added for later will become the top priority. Such is
the shifting world of ToDo lists. To accommodate this Taskforge has the ``task
workon`` command. To demonstrate let's make ``another default priority task the
top priority``. To do this let's find its ID with ``task todo``:

.. console::

   $ task todo
   | ID                               | Created Date               | Completed Date | Priority | Title                           | Context |
   | -------------------------------- | -------------------------- | -------------- | -------- | ------------------------------- | ------- |
   | eabdeee413ef442fa68c994119d817d2 | 2018-09-23 18:41:18.858741 | None           | 1.0      | complete the taskforge tutorial | default |
   | 1e634ced06d64093a747f38da024f9a6 | 2018-09-23 18:46:05.198426 | None           | 1.0      | another default priority task   | default |
   $

Then run the ``task workon`` command providing the ID of the task we want to
re-prioritize:

.. console::

   $ task workon 1e634ced06d64093a747f38da024f9a6
   $


``task next`` should now show ``another default priority task`` as the
current task:

.. console::

   $ task next
   1e634ced06d64093a747f38da024f9a6: another default priority task
   $

It accomplishes this by determining the priority of the current task and adding
``0.1`` to it. If we run ``task todo`` we can see this:

.. console::

   $ task todo
   | ID                               | Created Date               | Completed Date | Priority | Title                           | Context |
   | -------------------------------- | -------------------------- | -------------- | -------- | ------------------------------- | ------- |
   | eabdeee413ef442fa68c994119d817d2 | 2018-09-23 18:41:18.858741 | None           | 1.0      | complete the taskforge tutorial | default |
   | 1e634ced06d64093a747f38da024f9a6 | 2018-09-23 18:46:05.198426 | None           | 1.1      | another default priority task   | default |
   $

Let's go ahead and complete this task now. A shortcut that we did not mention
earlier is that if ``task done`` is given no arguments it will complete the
current task:

.. console::

   $ task done
   $ task next
   eabdeee413ef442fa68c994119d817d2: complete the taskforge tutorial
   $

This is a useful shortcut since most often you'll be completing the current task
as you work through your task list.

Further Reading
---------------

You can safely run ``task done`` now since you've completed the getting started
guide for Taskforge. From here you can start looking at using different
:doc:`../lists/index` or see the :doc:`../advanced_usage/index` guide to find out how
to integrate Taskforge with external reporting tools.

- :doc:`../tutorials/index`
- :doc:`../how_tos/configuring_taskforge`
- :doc:`../how_tos/query_language`
- :doc:`../lists/index`
- :doc:`../advanced_usage/index`


