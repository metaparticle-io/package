#!/usr/bin/env python
'''Unit tests for DockerRunner'''


import unittest
from mock import patch, call, MagicMock
from metaparticle_pkg.runner import docker


class TestDockerRunner(unittest.TestCase):
    '''Unit tests for DockerRunner'''

    def setUp(self):
        self.docker_runner = docker.DockerRunner()

    @patch("metaparticle_pkg.runner.docker.subprocess")
    def test_run(self, mocked_subprocess):
        '''Test Run method'''

        # Input arguments
        img = "test_image"
        name = "test_name"
        options = MagicMock()
        options.ports = ["4562"]

        # Expected argument called with
        expected_args = [
            'docker', 'run', '-d',
            '--name', 'test_name',
            '-p', '4562:4562', 'test_image'
        ]

        self.docker_runner.run(img, name, options)

        mocked_subprocess.check_call.assert_called_once_with(
            expected_args)

    @patch("metaparticle_pkg.runner.docker.subprocess")
    def test_logs(self, mocked_subprocess):
        '''Test logs method'''

        # Input arguments
        name = "test_name"

        # Expected argument called with
        expected_args = ['docker', 'logs', '-f', name]

        self.docker_runner.logs(name)

        mocked_subprocess.check_call.assert_called_once_with(
            expected_args)

    @patch("metaparticle_pkg.runner.docker.subprocess")
    def test_cancel(self, mocked_subprocess):
        '''Test cancel method'''

        # Input arguments
        name = "test_name"

        # Expected calls with supplied arguments
        expected_call_list = [
            call(['docker', 'kill', name]),
            call(['docker', 'rm', name])
        ]

        self.docker_runner.cancel(name)

        mocked_subprocess.check_call.assert_has_calls(
            expected_call_list)


if __name__ == '__main__':
    unittest.main()
