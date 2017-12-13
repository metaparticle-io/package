# Metaparticle for Python
Metaparticle/Package is a collection of libraries intended to
make building and deploying containers a seamless and idiomatic
experience for developers.

This is the implementation for Python (3.x, tested on 3.6.x).

## Introduction
Metaparticle/Package simplifies and centralizes the task of
building and deploying a container image.

Here is a quick example.

Consider this simple Python application:
```python
def main():
    print('Hello World')


if __name__ == '__main__':
   main()
```

To containerize this application, you need to use the `metaparticle` library and
the `containerize` wrapper function like this:

```python
from metaparticle import Containerize


@Containerize(package={'name': 'testcontainer', 'repo': 'brendanburns', 'publish': True})
def main():
    print('hello world')


if __name__ == '__main__':
    main()

```

When you run this application, instead of printing "Hello world", it first packages itself as a container, and
then (optionally) deploys itself inside that container.

## Tutorial

```bash
git clone https://github.com/metaparticle-io/package/
cd package/python

make venv
source venv/bin/activate

cd examples
python example.py
```