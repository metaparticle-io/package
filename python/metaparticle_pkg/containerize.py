import os
import signal
import sys 

import metaparticle_pkg.option as option
import metaparticle_pkg.builder as builder
import metaparticle_pkg.runner as runner

def is_in_docker_container():
    mp_in_container = os.getenv('METAPARTICLE_IN_CONTAINER', None)
    if mp_in_container in ['true', '1']:
        return True
    elif mp_in_container in ['false', '0']:
        return False

    try:
        with open('/proc/1/sched', 'rt') as f:
            if '(1,' in f.readline():
                return False
            else:
                return True

    except FileNotFoundError:
        return False


def write_dockerfile(package):
    with open('Dockerfile', 'w+t') as f:
        f.write("""FROM python:{version}-alpine

COPY ./ /{name}/
RUN pip install -r /{name}/requirements.txt

CMD python /{name}/example.py
""".format(name=package.name, version=package.py_version))


class Containerize(object):

    def __init__(self, runtime={}, package={}):
        self.runtime = option.load(option.RuntimeOptions, runtime)
        self.package = option.load(option.PackageOptions, package)
        self.image = "{repo}/{name}:latest".format(repo=self.package.repository, name=self.package.name)

        self.builder = builder.select(self.package.builder)
        self.runner = runner.select(self.runtime.executor)

    def __call__(self, func):
        def wrapped(*args, **kwargs):
            if is_in_docker_container():
                return func(*args, **kwargs)

            write_dockerfile(self.package)
            self.builder.build(self.image)

            if self.package.publish:
                self.builder.publish(self.image)

            def signal_handler(signal, frame):
                self.runner.cancel(name)
                sys.exit(0)
            signal.signal(signal.SIGINT, signal_handler)

            self.runner.run(self.image, self.package.name, self.runtime)
            return self.runner.logs(self.image)
        return wrapped