from __future__ import absolute_import
import logging

from docker import APIClient

# use a generic logger name: metaparticle_pkg.runner
logger = logging.getLogger('.'.join(__name__.split('.')[:-1]))


class DockerRunner:
    def __init__(self):
        self.docker_client = None

    def run(self, img, name, options):
        if self.docker_client is None:
            self.docker_client = APIClient(version='auto')

        ports = []
        host_config = None

        # Prepare port configuration
        if options.ports is not None and len(options.ports) > 0:
            for port_number in options.ports:
                ports.append(port_number)

            host_config = self.docker_client.create_host_config(
                port_bindings={p: p for p in ports}
            )

        # Launch docker container
        container = self.docker_client.create_container(
            img,
            host_config=host_config,
            name=name,
            ports=ports
        )

        logger.info('Starting container {}'.format(container))

        self.docker_client.start(container=container.get('Id'))
        self.container = container

    def logs(self, *args, **kwargs):
        if self.docker_client is None:
            self.docker_client = APIClient(version='auto')

        log_stream = self.docker_client.logs(
            self.container.get('Id'),
            stream=True,
            follow=True
        )

        for line in log_stream:
            logger.info(line.decode("utf-8").strip('\n'))

    def cancel(self, name):
        if self.docker_client is None:
            self.docker_client = APIClient(version='auto')
        self.docker_client.kill(self.container.get('Id'))
        self.docker_client.remove_container(self.container.get('Id'))
