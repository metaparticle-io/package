from metaparticle_pkg.builder.docker import DockerBuilder
def select(spec):
    if spec == 'docker':
        return DockerBuilder()
    raise Exception('Unknown spec {}'.format(spec))