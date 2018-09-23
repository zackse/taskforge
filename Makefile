PYTHON := python3

lint:
	$(PYTHON) -m pydocstyle taskforge
	$(PYTHON) -m pylint taskforge tests

fmt:
	$(PYTHON) -m yapf --recursive -i taskforge tests

clean:
	rm -rf build dist
	find . -regex '.*egg-info' -type d -exec rm -rf {} \;
	find . -name '*.pyc'

install:
	$(PYTHON) setup.py install

install-dev:
	pip install --editable .
	pip install --editable ".[mongo]"
	pip install -r requirements.dev.txt

test:
	PYTHONPATH="$$PYTHONPATH:src" $(PYTHON) -m pytest -m 'not slow'

test-all:
	PYTHONPATH="$$PYTHONPATH:src" $(PYTHON) -m pytest --disable-pytest-warnings

.PHONY: docs
docs: html
	mv src/docs/build/html/* docs/

# You can set these variables from the command line.
SPHINXOPTS    =
SPHINXBUILD   = sphinx-build
SOURCEDIR     = src/docs
BUILDDIR      = build

# Put it first so that "make" without argument is like "make help".
help:
	@$(SPHINXBUILD) -M help "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)

.PHONY: help Makefile

livehtml:
	sphinx-autobuild --watch ../taskforge -b html $(SPHINXOPTS) "$(SOURCEDIR)" $(BUILDDIR)/html

# Catch-all target: route all unknown targets to Sphinx using the new
# "make mode" option.  $(O) is meant as a shortcut for $(SPHINXOPTS).
%: Makefile
	@$(SPHINXBUILD) -M $@ "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)
