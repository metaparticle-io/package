#!/usr/bin/env python
'''Unit tests for DockerBuilder'''


import unittest
from mock import patch
from metaparticle_pkg.builder import docker


class TestDockerBuilder(unittest.TestCase):
    '''Unit tests for DockerBuilder'''

    def setUp(self):
        self.docker_builder = docker.DockerBuilder()

        # Input arguments
        self.img = "test_image"

    @patch("metaparticle_pkg.builder.docker.subprocess")
    def test_build(self, mocked_subprocess):
        '''Test build method'''

        # Expected argument called with
        expected_args = ['docker', 'build', '-t', self.img, '.']

        self.docker_builder.build(self.img)

        mocked_subprocess.check_call.assert_called_once_with(
            expected_args)

    @patch("metaparticle_pkg.builder.docker.subprocess")
    def test_publish(self, mocked_subprocess):
        '''Test publish method'''

        # Expected argument called with
        expected_args = ['docker', 'push', self.img]

        self.docker_builder.publish(self.img)

        mocked_subprocess.check_call.assert_called_once_with(
            expected_args)


if __name__ == '__main__':
    unittest.main()
