import functools
import os


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

    except FileNotFoundError:
        return False


def write_dockerfile(name):
    with open('Dockerfile', 'w+t') as f:
        f.write("""FROM python:3

COPY ./ /{}/
RUN pip install -r /{}/requirements.txt

CMD python /{}/main.py
    """.format(name, name, name))


def select_builder(buildSpec):
    if buildSpec == 'docker':
        return __import__('docker_builder')
    else:
        raise Exception('Unknown buildSpec {}'.format(buildSpec))


def select_runner(runSpec):
    if runSpec == 'docker':
        return __import__('docker_runner')
    else:
        raise Exception('Unknown buildSpec {}'.format(buildSpec))


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
                return
            else:
                name = get_name(options)
                img = "{}/{}:latest".format(repository, name)
                builder = select_builder('docker')
                runner = select_runner('docker')

            write_dockerfile(name)
            builder.build(img)

            try:
                if options['publish']:
                    builder.publish(img)

            except KeyError:
                pass

            runner.run(img, name, options)

            return f(*args, **kwargs)
        return wrapper
    return _containerize
