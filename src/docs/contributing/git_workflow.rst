Git Workflow
============

Taskforge uses a "no merge commits" strategy for managing git branches and pull
requests.

Whenever you are about to submit a PR follow these steps to ensure a smooth
review process. Make sure to replace ``origin`` in these examples with whatever
remote you have the primary repository configured as. At the time of this
writing that would be ``http://github.com/chasinglogic/taskforge``.

1. Download the latest changes using ``git fetch --all``
2. Now that you've downloaded the latest changes you can "rebase" your working
   branch on the latest master. A detailed description of rebasing is outside
   the scope of this article, but in english it will re-apply your changes on
   top of master. To do this run:

.. code::

   git rebase origin/master

3. Now your changes are the most recent commits. You'll need to push your
   rebased branch to the remote, if you've already pushed before you'll need to
   do a force push. The command for doing this is:

.. code::

   git push --force-with-lease

4. Now that's done you can submit your PR and your branch will have the latest
   changes from master.

Git Best Practices
==================

Here are some best practices that while not strictly enforced, can come up
during code review:

Commit Message Formatting
+++++++++++++++++++++++++

All commit messages should ideally follow this format:

.. code::

   A commit summary that is less than 80 characters.

   Followed by a blank newline and a longer description if necesary. Not all
   changes will need a description like this one. The 80 character limit in the
   summary is not a hard rule so use your best judgement.

   If you're combining more than one fix or chang into a commit you should list
   them in a bullet point format like so:

    - Fix typo in docstring
    - Corrected bad link in README

With the following format, it's acceptable to throw smaller changes into one
commit with a bigger change. The most important thing is that I should be able
to tell what the commit changed by reading your message and not the diff.

The added benefit of this format is that GitHub will automatically populate your
Pull Request with this information.

We will not merge commit messages that are offensive, i.e. violating our
:doc:`code_of_conduct`, or non-helpful, such as ``remove incorrectly added
file``. You will need to rebase and drop these commits and correct your commit
history first.

No Merge Commits
++++++++++++++++

We avoid merge commits on master, so if your pull request includes them we will
ask you to remove them.


Commit Limits
+++++++++++++


While for some features it's understandable that you will have a lot of commits,
generally speaking, a good Pull Request should not generally contain any more
than 5 commits. Again this is not a hard and fast rule, but if you find yourself
with 5+ commits rethink how you've written your commits and consider squashing
some of them.

Topic Branches
++++++++++++++

When working on a bug or feature you should create a topic branch for that work.
This isolates your commits in a way that prevents you from having to remove
erroneous commits later. A topic branch should always come from master and
should have a descriptive name, at least to you. The commands for creating a
topic branch are:

.. code::

   git fetch --all
   git branch origin/master my-topic-branch
   git checkout my-topic-branch

This does make the assumption that ``origin`` is the master repository. If
``origin`` is instead pointing to your fork change it to whatever your forks
remote name is.

Further Resources
=================

Here are some valuable resources you can use when you find yourself stuck on
git or just want a better understanding of the topics mentioned here:

- `Writing a good commit message <https://chris.beams.io/posts/git-commit/>`_
- `Introduction to working with git history <https://robots.thoughtbot.com/git-interactive-rebase-squash-amend-rewriting-history>`_
- `Git communities <https://git-scm.com/community>`_
- `Rebase reference page <https://git-scm.com/docs/git-rebase>`_
