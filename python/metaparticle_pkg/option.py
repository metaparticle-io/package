from collections import namedtuple
import sys
import os

def load(cls, options):
    if not isinstance(options, dict):
        sys.stderr.write("Must provide an options dict.")
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

class RuntimeOptions(namedtuple('Runtime', 'executor replicas ports publicAddress shardSpec')):
    required_options = []
    def __new__(cls, executor='docker', replicas=1, ports=[], publicAddress=[], shardSpec=None):
        return super(RuntimeOptions, cls).__new__(cls, executor, replicas, ports, publicAddress, shardSpec)

class ShardSpec(namedtuple('ShardSpec', 'shards shardExpression')):
    required_options = []
    def __new__(cls, shards=0, shardExpression='.*'):
        return super(ShardSpec, cls).__new__(cls, shards, shardExpression)


class PackageOptions(namedtuple('Package', 'repository name builder publish verbose quiet py_version')):
    required_options = ['repository']
    def __new__(cls, repository, name, builder='docker', publish=False, verbose=True, quiet=False, py_version=3, dockerfile=None):
        name = name if name else os.path.basename(os.getcwd())
        return super(PackageOptions, cls).__new__(cls, repository, name, builder, publish, verbose, quiet, py_version)
