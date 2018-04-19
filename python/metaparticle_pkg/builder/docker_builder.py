import json
import logging

from docker import APIClient

logger = logging.getLogger('.'.join(__name__.split('.')[:-1]))


class DockerBuilder:
    def __init__(self):
        self.docker_client = None

    def build(self, img, path='.'):
        if self.docker_client is None:
            self.docker_client = APIClient(version='auto')

        bld = self.docker_client.build(
            path=path,
            tag=img,
            encoding='utf-8'
        )

        for line in bld:
            self._process_stream(line)

    def publish(self, img):
        if self.docker_client is None:
            self.docker_client = APIClient(version='auto')

        # TODO: do we need to set tag?
        for line in self.docker_client.push(img, stream=True):
            self._process_stream(line)

    def _process_stream(self, line):
        ln = line.decode('utf-8').strip()

        # try to decode json
        try:
            ljson = json.loads(ln)

            if ljson.get('error'):
                msg = str(ljson.get('error', ljson))
                logger.error('Build failed: ' + msg)
                raise Exception('Image build failed: ' + msg)
            else:
                if ljson.get('stream'):
                    msg = 'Build output: {}'.format(ljson['stream'].strip())
                elif ljson.get('status'):
                    msg = 'Push output: {} {}'.format(
                        ljson['status'],
                        ljson.get('progress')
                    )
                elif ljson.get('aux'):
                    msg = 'Push finished: {}'.format(ljson.get('aux'))
                else:
                    msg = str(ljson)

                logger.debug(msg)

        except json.JSONDecodeError:
            logger.warning('JSON decode error: {}'.format(ln))
