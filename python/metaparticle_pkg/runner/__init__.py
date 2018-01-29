from metaparticle_pkg.runner.docker_runner import DockerRunner
from metaparticle_pkg.runner.metaparticle import MetaparticleRunner


def select(spec):
    if spec == 'docker':
        return DockerRunner()
    if spec == 'metaparticle':
        return MetaparticleRunner()

    raise Exception('Unknown spec {}'.format(spec))
