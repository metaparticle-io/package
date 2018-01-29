from metaparticle_pkg.builder.docker_builder import DockerBuilder


def select(spec):
    if spec == 'docker':
        return DockerBuilder()
    raise Exception('Unknown spec {}'.format(spec))
