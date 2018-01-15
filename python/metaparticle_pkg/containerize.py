import os
import shutil
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
        with open('/proc/1/cgroup', 'r+t') as f:
            lines = f.read().splitlines()
            last_line = lines[-1]
            if 'docker' in last_line:
                return True
            elif 'kubepods' in last_line:
                return True
            else:
                return False

    except IOError:
        return False


def write_dockerfile(package, exec_file):
    if hasattr(package, 'dockerfile') and package.dockerfile is not None:
        shutil.copy(package.dockerfile, 'Dockerfile')
        return

    with open('Dockerfile', 'w+t') as f:
        f.write("""FROM python:{version}-alpine

COPY ./ /app/
RUN pip install --no-cache -r /app/requirements.txt

CMD python /app/{exec_file}
""".format(version=package.py_version, exec_file=exec_file))


class Containerize(object):

    def __init__(self, runtime={}, package={}):
        self.runtime = option.load(option.RuntimeOptions, runtime)
        self.package = option.load(option.PackageOptions, package)
        self.image = "{repo}/{name}:latest".format(
            repo=self.package.repository,
            name=self.package.name
        )

        self.builder = builder.select(self.package.builder)
        self.runner = runner.select(self.runtime.executor)

    def __call__(self, func):
        def wrapped(*args, **kwargs):
            if is_in_docker_container():
                return func(*args, **kwargs)

            exec_file = sys.argv[0]
            slash_ix = exec_file.find('/')
            if slash_ix != -1:
                exec_file = exec_file[slash_ix:]

            write_dockerfile(self.package, exec_file)
            self.builder.build(self.image)

            if self.package.publish:
                self.builder.publish(self.image)

            def signal_handler(signal, frame):
                self.runner.cancel(self.package.name)
                sys.exit(0)
            signal.signal(signal.SIGINT, signal_handler)

            self.runner.run(self.image, self.package.name, self.runtime)

            return self.runner.logs(self.package.name)
        return wrapped
