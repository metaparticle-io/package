#!/usr/bin/python
from metaparticle_pkg import Containerize

import time
import logging

# all metaparticle output is accessible through the stdlib logger (debug level)
logging.basicConfig(level=logging.INFO)
logging.getLogger('metaparticle_pkg.runner').setLevel(logging.DEBUG)
logging.getLogger('metaparticle_pkg.builder').setLevel(logging.DEBUG)


@Containerize(
    package={
        'name': 'simple',
        'repository': 'docker.io/brendanburns',
        'publish': False
    }
)
def main():
    print('hello world!')

    for i in range(5):
        print('Sleeping ... {} sec'.format(i))
        time.sleep(1)


if __name__ == '__main__':
    main()
