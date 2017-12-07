# Metaparticle for Python
Metaparticle/Package is a collection of libraries intended to
make building and deploying containers a seamless and idiomatic
experience for developers.

This is the implementation for Python.

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
from metaparticle import containerize


@containerize('brendanburns', options={'name': 'testcontainer', 'publish': True})
def main():
    print('hello world')


if __name__ == '__main__':
    main()

```

When you run this application, instead of printing "Hello world", it first packages itself as a container, and
then (optionally) deploys itself inside that container.
