PYTHON := python3

lint:
	$(PYTHON) -m pydocstyle taskforge
	$(PYTHON) -m pylint taskforge tests

fmt:
	$(PYTHON) -m yapf --recursive -i taskforge tests

clean:
	rm -rf build dist
	find . -regex '.*egg-info' -type d -exec rm -rf {} \;
	find . -path ./.venv -prune -type f -name '*.pyc'

build-docs:
	$(PYTHON) src/docs/build_docs.py

install:
	$(PYTHON) setup.py install

install-dev:
	pip install --editable .
	pip install --editable ".[mongo]"
	pip install -r requirements.dev.txt

test:
	PYTHONPATH="$$PYTHONPATH:src" $(PYTHON) -m pytest
