import setuptools
import json

from os import path

base_path = path.abspath(path.dirname(__file__))
version_json = path.join(base_path, 'metaparticle_pkg/version.json')

with open(version_json) as f:
    config = json.loads(f.read())

setuptools.setup(
    name='metaparticle_pkg',
    version=config['version'],
    url='https://github.com/metaparticle-io/package/tree/master/python',
    license=config['license'],
    description='Easily containerize your python application',
    author='Metaparticle Authors',
    packages=setuptools.find_packages(),
    package_data={},
    include_package_data=False,
    zip_safe=False,
    install_requires=['docker==2.7.0'],
    test_require=['pytest', 'flake8'],
    platforms='linux',
    keywords=['kubernetes', 'docker', 'container', 'metaparticle'],
    # latest from https://pypi.python.org/pypi?%3Aaction=list_classifiers
    classifiers=[
        'Development Status :: 4 - Beta',
        'Environment :: Console',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: MIT License',
        'Natural Language :: English',
        'Operating System :: POSIX :: Linux',
        'Programming Language :: Python',
        'Programming Language :: Python :: 2',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: Implementation :: CPython',
        'Topic :: Software Development :: Libraries :: Python Modules',
        'Topic :: Utilities',
        ]
)
