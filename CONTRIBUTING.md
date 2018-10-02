# Contributing

The basic workflow for contributing looks like this:

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. :fire: Submit a pull request :D :fire:

All pull requests should go to the master branch. Thanks!

## Contributing Code

For setting up a local dev environment you will need Python 3 and a virtualenv.

From there you can run `make install-dev` this will install development
dependencies as well as tools. Additionally, it will install a "editable"
version of the task command that will allow you to test changes without
reloading or reinstalling after a change.

Your code will need tests and must pass `make test`. Additionally no errors
should show up with `make lint` which uses pylint and pydocstyle to enforce
a consistent style guide.

If you have suggestions for changes to the style guide file an issue before
making your changes.


## Contributing Docs

Currently I am working on a Vale style guide and integrating this linter into
our CI system. As it stands you can view your documentation changes with 
`make livehtml` which will run a local html server that live reloads changes to
the docs.

All docs are kept in `src/docs` and written in restructured text. They are build
using Sphinx.
