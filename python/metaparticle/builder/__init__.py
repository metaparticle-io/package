import importlib

def select(spec):
    if spec == 'docker':
        return importlib.import_module('.docker', 'metaparticle.builder')
    else:
        raise Exception('Unknown spec {}'.format(spec))