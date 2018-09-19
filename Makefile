lint:
	pydocstyle taskforge
	pylint taskforge tests

fmt:
	yapf --recursive -i taskforge tests

clean:
	rm -rf build dist
	find . -regex '.*egg-info' -type d -exec rm -rf {} \;
	find . -path ./.venv -prune -type f -name '*.pyc'

install:
	python setup.py install

install-dev:
	pip install --editable .
	pip install --editable ".[mongo]"
	pip install yapf pydocstyle pylint

test:
	PYTHONPATH="$$PYTHONPATH:src" python -m unittest discover
