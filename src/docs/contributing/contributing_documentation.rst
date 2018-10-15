Contributing Documentation
==========================

Contributing documentation to Taskforge, like all contributions, is highly
appreciated. This document will help you set up your environment and follow our
documentation best practices.

While not required for writing documentation if you would like to build the docs
and see your changes as they'll be on the website you will need Python version
3.4 or greater installed. If you need to install Python, you can `go to Python's
website <https://python.org>`_ to install it on your platform.

Getting the Code
++++++++++++++++

.. note::

   GitHub does provide a web editor and automates some of this process, if you
   just want to make edits to an existing document you can skip this section and
   click the "pencil" icon in the top right of the document you want to edit.

Before you begin writing documentation for Taskforge you will need to download
the repository. To do this open your shell of choice and run the following
command:

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

Requirements for Submitting Documentation
+++++++++++++++++++++++++++++++++++++++++

All documentation needs to meet these requirements:

- Passes our documentation lint step in CI
- Uses no offensive language
- Does not break rendering of the HTML site
- Must live under ``src/docs`` and written in `reStructuredText
  <http://www.sphinx-doc.org/en/master/usage/restructuredtext/basics.html>`_

We lint our documentation using `Vale <https://github.com/errata-ai/vale>`_ it
is a CLI tool for linting prose and is syntax aware. You can install it using
their instructions at the link above.

It's not required that you run the linting locally but will save you time since
Travis CI can take some to return a result. You can lint your documents by
passing them to Vale as shown:

.. console::

   $ vale src/docs/$YOUR_DOCUMENT_PATH

How We Categorize Documentation
+++++++++++++++++++++++++++++++

The Taskforge project has 4 categories of documentation. How-to's, Tutorials,
Designs, and Usage / Reference documentation. Each category is explained below:

- A tutorial:

  - is learning-oriented
  - allows the newcomer to get started
  - is a lesson

Analogy: teaching a small child how to cook
How-to guides

- A how-to guide:

  - is goal-oriented
  - shows how to solve a specific problem
  - is a series of steps

Analogy: a recipe in a cookery book
Explanation

- A design:

  - is understanding-oriented
  - explains
  - provides background and context

Analogy: an article on culinary social history
Reference

- A usage guide:

  - is information-oriented
  - describes the machinery
  - is accurate and complete

Analogy: a reference encyclopaedia article

These categories and explanations are taken from `this article
<https://www.divio.com/blog/documentation/>`_ by Daniele Procida.

Most of these have corresponding folders in ``src/docs``. Try to place your
documents in the appropriate category folder using your best judgement. If
you're not sure pick one and it can be discussed during code review.
