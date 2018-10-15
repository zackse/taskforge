Contributing Code
=================

Contributing code to Taskforge, like all contributions, is highly appreciated.
This document will help you set up your development environment and follow our
coding best practices.

This document assumes that you have already installed Python version 3.4 or
greater. If not you can `go to Python's website <https://python.org>`_ to
install it on your platform.

.. contents::

.. note::

   Anywhere you see ``python`` if you're on Mac OS or an older Linux distro then
   you will need to change it to ``python3``.

Getting the Code
++++++++++++++++

Before you begin development on Taskforge you will need to download the
repository. To do this open your shell of choice and run the following command:

.. console::

   $ git clone https://github.com/chasinglogic/taskforge

Now use GitHub to fork the project to create your own remote which you will work
from. When viewing the repo above in your browser, click the fork button in the
top right hand corner.

GitHub provides `this article <https://help.github.com/articles/fork-a-repo/>`_
which gives a good explanation of what forking means and how to work with a fork.

Once you've created your fork copy the clone URL and run the command:

.. console::

   $ git remote add fork $YOUR_CLONE_URL

Replacing ``$YOUR_CLONE_URL`` with the URL you just copied from GitHub.

Setting Up A Development Environment
++++++++++++++++++++++++++++++++++++

.. note::

   It's recommended that you set up a Python virtualenv before doing any of
   the steps below. A great explanation of how and why to use a virtualenv can
   be found at `https://docs.python.org/3/library/venv.html`_

First lets cd into the repository we created earlier if you haven't already:

.. console::

   $ cd taskforge

To install the development tools and libraries you need you can run these
commands:

.. console::

   $ python -m pip install --editable .
   $ python -m pip install -r requirements.dev.txt

This will install the testing and linting tools you'll need to make sure your
code is ready for review.

Requirements for Submitting Code
++++++++++++++++++++++++++++++++

All code needs to meet these requirements:

- If new code:
  - Write a design document and complete the :doc:`processes/design`_ 
  - For each goal in the design write a test
  - Write the code to make the tests pass in CI
- If fixing a bug:
  - Write a test which reproduces the bug
  - Write the code fixing that test, it must pass in CI
- All code must pass lint using our pylintrc which is in the root of the
  repository

You can run the linting steps locally with these commands. Although, it's worth
noting most text editors will integrate with these tools automatically:

.. console::

   $ python -m pylint src tests
   $ python -m pydocstyle src

For testing we use pytest. To run the test suite you can use the command:

.. console::

   $ PYTHONPATH="$PYTHONPATH:src" python -m pytest -m 'not slow'

Any tests which call external services or databases must have the pytest marker
indicating it as slow. To run those tests, remove the marker flag from the
previous command:

.. console::

   $ PYTHONPATH="$PYTHONPATH:src" python -m pytest

.. note::

   For unix systems which have ``make`` installed you can perform the above
   commands with:

   .. code::

      $ make lint
      $ make test
      $ make test-all
