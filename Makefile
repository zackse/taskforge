lint:
	pydocstyle taskforge
	pylint taskforge tests

fmt:
	yapf --recursive -i taskforge tests

clean:
	rm -rf *.egg-info build dist
	find . -path ./.venv -prune -type f -name '*.pyc'

install:
	python setup.py install

install-dev:
	pip install --editable .
	pip install -e --editable .[cli]
	pip install -r requirements.dev.txt

test:
	python -m unittest discover
