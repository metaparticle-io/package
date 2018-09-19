#!/usr/bin/python
from metaparticle_pkg import Containerize, PackageFile

import os
import time
import logging

# all metaparticle output is accessible through the stdlib logger (debug level)
logging.basicConfig(level=logging.INFO)
logging.getLogger('metaparticle_pkg.runner').setLevel(logging.DEBUG)
logging.getLogger('metaparticle_pkg.builder').setLevel(logging.DEBUG)


DATA_FILE = '/opt/some/random/spot/data1.json'
SCRIPT = '/opt/another/random/place/get_the_data.sh'


@Containerize(
    package={
        'name': 'file-example',
        'repository': 'docker.io/brendanburns',
        'publish': False,
        'additionalFiles': [
            PackageFile(src='./data.json', dest=DATA_FILE, mode='0400'),
            PackageFile(src='./get_data.sh', dest=SCRIPT),
        ]
    }
)
def main():
    os.system(SCRIPT)
    for i in range(5):
        print('Sleeping ... {} sec'.format(i))
        time.sleep(1)


if __name__ == '__main__':
    main()
