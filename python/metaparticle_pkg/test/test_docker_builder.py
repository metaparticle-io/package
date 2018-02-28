#!/usr/bin/env python
'''Unit tests for DockerBuilder'''

import sys
import unittest
from metaparticle_pkg.builder import docker_builder

if sys.version_info >= (3, 3):
    from unittest.mock import patch
else:
    from mock import patch


class TestDockerBuilder(unittest.TestCase):
    '''Unit tests for DockerBuilder'''

    def setUp(self):
        self.builder = docker_builder.DockerBuilder()

        # Input arguments
        self.img = 'test_image'

    @patch('metaparticle_pkg.builder.docker_builder.APIClient.build')
    def test_build(self, mocked_build):
        '''Test build method'''

        self.builder.build(self.img)

        mocked_build.assert_called_once_with(
            path='.',
            tag=self.img,
            encoding='utf-8'
        )

    @patch('metaparticle_pkg.builder.docker_builder.APIClient.build',
           return_value=[b'{"error":"test"}', ])
    def test_build_with_docker_error(self, mocked_build):
        '''Test build method with error returned building image'''

        with self.assertRaises(Exception) as context:
            self.builder.build(self.img)

        self.assertTrue('Image build failed:' in str(context.exception))

    @patch('metaparticle_pkg.builder.docker_builder.APIClient.push')
    def test_publish(self, mocked_push):
        '''Test publish method'''

        # Expected argument called with
        self.builder.publish(self.img)

        mocked_push.assert_called_once_with(
            self.img,
            stream=True
        )

    @patch('metaparticle_pkg.builder.docker_builder.APIClient.push',
           return_value=[b'{"error":"test"}', ])
    def test_publish_with_docker_error(self, mocked_build):
        '''Test build method with error returned building image'''

        with self.assertRaises(Exception) as context:
            self.builder.publish(self.img)

        self.assertTrue('Image build failed:' in str(context.exception))


if __name__ == '__main__':
    unittest.main()
