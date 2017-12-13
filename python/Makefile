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

clean:
	rm -rf build dist *.egg-info MANIFEST