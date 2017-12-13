import importlib


def select(spec):
    if spec == 'docker':
        return importlib.import_module('.docker', 'metaparticle.runner')
    else:
        raise Exception('Unknown spec {}'.format(spec))