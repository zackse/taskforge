lint:
	pydocstyle taskforge
	pylint taskforge tests

fmt:
	yapf --recursive -i taskforge tests

clean:
	rm -rf build dist
	find . -regex '.*egg-info' -type d -exec rm -rf {} \;
	find . -path ./.venv -prune -type f -name '*.pyc'

build-docs:
	python src/docs/build_docs.py

install:
	python setup.py install

install-dev:
	pip install --editable .
	pip install --editable ".[mongo]"
	pip install -r requirements.dev.txt

test:
	PYTHONPATH="$$PYTHONPATH:src" python -m unittest discover
