import logging

from docker import APIClient

logger = logging.getLogger(__name__)


class DockerRunner:
    def __init__(self):
        self.docker_client = APIClient(version='auto')

    def run(self, img, name, options):
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
            name=name,
            ports=ports,
            host_config=host_config,
        )
        self.docker_client.start(container=container.get('Id'))

        self.container = container

        logger.info('Starting container {}'.format(container))

    def logs(self, *args, **kwargs):
        # seems like we are hitting bug
        # https://github.com/docker/docker-py/issues/300
        log_stream = self.docker_client.logs(
            self.container.get('Id'),
            stream=True,
            follow=True
        )

        for line in log_stream:
            logger.info(line)

    def cancel(self, name):
        self.docker_client.kill(self.container.get('Id'))
        self.docker_client.remove_container(self.container.get('Id'))
