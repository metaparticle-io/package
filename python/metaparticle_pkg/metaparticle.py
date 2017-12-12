import functools
import os
import signal
import sys
from metaparticle_pkg.docker_builder import DockerBuilder
from metaparticle_pkg.docker_runner import DockerRunner
from metaparticle_pkg.metaparticle_runner import MetaparticleRunner

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


def write_dockerfile(name, exec_file):
    with open('Dockerfile', 'w+t') as f:
        f.write("""FROM python:3-alpine

COPY ./ /{}/
RUN pip install -r /{}/requirements.txt

CMD python /{}/{}
    """.format(name, name, name, exec_file))


def select_builder(buildSpec):
    if buildSpec == 'docker':
        return DockerBuilder()
    raise Exception('Unknown buildSpec {}'.format(buildSpec))


def select_runner(runSpec):
    if runSpec == 'docker':
        return DockerRunner()
    if runSpec == 'metaparticle':
        return MetaparticleRunner()
    raise Exception('Unknown buildSpec {}'.format(runSpec))

def get_name(options):
    try:
        return options['name']
    except KeyError:
        pass

    return os.path.basename(os.getcwd())


def containerize(repository, options={}):
    def _containerize(f):

        @functools.wraps(f)
        def wrapper(*args, **kwargs):
            """
            Simplifies and centralizes the task of
            building and deploying a container image.
            """
            if is_in_docker_container():
                return f(*args, **kwargs)

            exec_file = sys.argv[0]
            slash_ix = exec_file.find('/')
            if slash_ix != -1:
                exec_file = exec_file[slash_ix:]

            name = get_name(options)
            img = "{}/{}:latest".format(repository, name)

            buildSpec = options.get('builder', 'docker')
            runSpec = options.get('runner', 'docker')
            builder = select_builder(buildSpec)
            runner = select_runner(runSpec)

            write_dockerfile(name, exec_file)
            builder.build(img)

            if options.get('publish', False):
                builder.publish(img)

            def signal_handler(signal, frame):
                print('Cancelling')
                runner.cancel(name)
                sys.exit(0)
            signal.signal(signal.SIGINT, signal_handler)

            runner.run(img, name, options)
            runner.logs(name)
        return wrapper
    return _containerize
