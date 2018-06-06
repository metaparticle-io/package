from __future__ import absolute_import
from collections import namedtuple
import sys
import os


def load(cls, options):
    if not isinstance(options, dict):
        sys.stderr.write("Must provide an options dictionary. Given option: %s" % repr(options))
        sys.exit(1)
    for option in cls.required_options:
        if option not in options:
            sys.stderr.write("Missing required field: %s" % option)
            sys.exit(1)
    try:
        return cls(**options)
    except TypeError as error:
        sys.stderr.write("Unexpected option(s) provided: %s" % error)
        sys.exit(1)


class RuntimeOptions(namedtuple('Runtime', 'executor replicas ports public shardSpec jobSpec')):
    required_options = []

    def __new__(cls, executor='docker', replicas=0, ports=[], public=False, shardSpec=None, jobSpec=None):
        if shardSpec:
            shardSpec = load(ShardSpec, shardSpec)
        if jobSpec:
            jobSpec = load(JobSpec, jobSpec)
        return super(RuntimeOptions, cls).__new__(cls, executor, replicas, ports, public, shardSpec, jobSpec)


class ShardSpec(namedtuple('ShardSpec', 'shards shardExpression')):
    required_options = []

    def __new__(cls, shards=0, shardExpression='.*'):
        return super(ShardSpec, cls).__new__(cls, shards, shardExpression)


class JobSpec(namedtuple('JobSpec', 'iterations')):
    required_options = ['iterations']

    def __new__(cls, iterations=0):
        return super(JobSpec, cls).__new__(cls, iterations)


class PackageOptions(namedtuple('Package', 'repository name builder publish py_version additionalFiles')):
    required_options = ['repository']

    def __new__(cls, repository, name, builder='docker', publish=False, py_version=3, dockerfile=None, additionalFiles=tuple()):
        name = name if name else os.path.basename(os.getcwd())
        return super(PackageOptions, cls).__new__(cls, repository, name, builder, publish, py_version, additionalFiles)
