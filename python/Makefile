.PHONY: test clean venv

test:
	(. venv/bin/activate; \
	python3 -m pytest --cov=. --cov-report=term-missing:skip-covered ;\
	python3 -m flake8 --config setup.cfg ;\
	)

venv: 
	virtualenv -p python3 venv
	venv/bin/pip install --upgrade pip
	venv/bin/pip install --upgrade setuptools
	. venv/bin/activate
	venv/bin/pip install -r requirements.txt
	venv/bin/pip install -e .

publish: 
	python setup.py sdist upload -r pypi

clean:
	rm -rf build dist *.egg-info MANIFEST
