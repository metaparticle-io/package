import setuptools
import os

exec(open('./metaparticle/version.py').read())

setuptools.setup(
    name='metaparticle',
    version=__version__,
    url='https://github.com/metaparticle-io/package/tree/master/python',
    license=__license__,
    description='Easily containerize your python application',
    packages=setuptools.find_packages(),
    package_data={},
    include_package_data=False,
    zip_safe=False,
    install_requires=[],
    platforms='linux',
    keywords=['kubernetes', 'docker', 'container', 'metaparticle'],
    # latest from https://pypi.python.org/pypi?%3Aaction=list_classifiers
    classifiers = [
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
        'Topic :: System :: Containers',
        'Topic :: Utilities',
        ]
)
