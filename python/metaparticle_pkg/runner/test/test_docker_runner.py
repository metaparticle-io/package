#!/usr/bin/env python
'''Unit tests for DockerRunner'''


import unittest
from unittest.mock import patch, MagicMock
from metaparticle_pkg.runner import docker_runner


class TestDockerRunner(unittest.TestCase):
    '''Unit tests for DockerRunner'''

    def setUp(self):
        self.runner = docker_runner.DockerRunner()

        # mock container
        self.cid = '123'
        self.runner.container = MagicMock()
        self.runner.container.get = MagicMock(return_value=self.cid)

    @patch('metaparticle_pkg.runner.docker_runner.APIClient.create_container')
    @patch('metaparticle_pkg.runner.docker_runner.APIClient.start')
    def test_run(self, mocked_start, mocked_create):
        '''Test Run method'''

        # Input arguments
        img = "test_image"
        name = "test_name"
        port = '4562'
        options = MagicMock()
        options.ports = [port]

        # Expected argument called with
        self.runner.run(img, name, options)

        mocked_create.assert_called_once_with(
            'test_image',
            host_config={
                'NetworkMode': 'default',
                'PortBindings': {
                    '{}/tcp'.format(port): [
                        {
                            'HostIp': '',
                            'HostPort': port,
                        }
                    ]
                }
            },
            name='test_name',
            ports=[port]
        )

        mocked_start.assert_called_once()

    @patch('metaparticle_pkg.runner.docker_runner.APIClient.logs')
    def test_logs(self, mocked_logs):
        '''Test logs method'''

        self.runner.logs()
        mocked_logs.assert_called_once_with('123', follow=True, stream=True)

    @patch('metaparticle_pkg.runner.docker_runner.APIClient.kill')
    @patch('metaparticle_pkg.runner.docker_runner.APIClient.remove_container')
    def test_cancel(self, mocked_remove, mocked_kill):
        '''Test cancel method'''

        # Input arguments
        name = "test_name"

        self.runner.cancel(name)

        mocked_kill.assert_called_with(self.cid)
        mocked_remove.assert_called_with(self.cid)


if __name__ == '__main__':
    unittest.main()
