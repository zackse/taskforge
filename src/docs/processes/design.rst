Design Document Process
=======================

This describes the design document process for Taskforge. All new features must
go through the design process as described here before the any code is written.

Why We Write Design Documents
-----------------------------

When the idea for a feature is first introduced we, as human beings, can look at
it with rose-colored glasses. It's much easier to just write code for features
than to think about them and, often once we start implementing them we realize
that the problem is both more difficult and has more edge cases than we
originally imagined.

Moving all ideas and features through a formal design process greatly improves
code quality and ensures that everything we build fits with the vision of the
maintainers of the project. Another benefit is that it gives us a reference point to look back on when wondering why a feature was created, or how that
feature was implemented.

The Design Process
------------------

The design process is as follows:

1. Find an issue labelled as "needs design"
2. As for any code or documentation change, clone the repository and make a
   branch. On this branch write a design doc then submit a Pull Request.

   - Design docs go under ``src/docs/designs``
   - All design docs should use the design doc template found at
     ``src/docs/designs/template.rst``
   - List designs have their own template which is slightly different. It is
     located at ``src/docs/designs/list_template.rst``

3. The maintainers will review your pull request
4. Once the maintainers have approved your design the pull request will be
   labelled with "Request for Comments"

   - From this point anyone in the Taskforge community have a chance to provide
     input and/or further the discussion.
   - If there is a disagreement the maintainers will have resolve them and
     have final say on changes to the design

5. The Pull Request will stay in the RFC phase for 10 days
6. After 10 days if all comments and issues have been resolved the PR will be
   merged and feature will be ready for work.
7. The maintainers will create a new issue to implement the design and close the
   "needs design" issue.

You are then free to pick up the implementation issue or leave it for someone
else to complete.

Use the :doc:`/designs/template` to get started.
