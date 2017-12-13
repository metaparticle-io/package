.PHONY: test clean venv

test:
	(. venv/bin/activate; \
	tox; \
	)

venv: 
	virtualenv -p python3 venv
	venv/bin/pip install --upgrade pip
	venv/bin/pip install --upgrade setuptools
	. venv/bin/activate
	venv/bin/pip install -e .
	venv/bin/pip install --upgrade tox

publish: 
	python setup.py sdist upload -r pypi

clean:
	rm -rf build dist *.egg-info MANIFEST
