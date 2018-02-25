#!/usr/bin/env python
'''Unit tests for containerize module'''

import unittest
from unittest.mock import patch, mock_open, MagicMock
from metaparticle_pkg import containerize
from types import FunctionType


class ContainerizeTest(unittest.TestCase):
    '''Unit tests for containerize module'''

    @patch("metaparticle_pkg.containerize.os")
    def test_is_in_docker_container_with_true_in_env(self, mocked_os):
        '''Test is_in_docker_container method for scenario
            when METAPARTICLE_IN_CONTAINER environment
            variable contains value 'true'.
        '''
        expected_value = True
        mocked_os.getenv.return_value = 'true'

        actual_value = containerize.is_in_docker_container()

        self.assertEqual(actual_value, expected_value)

    @patch("metaparticle_pkg.containerize.os")
    def test_is_in_docker_container_with_one_in_env(self, mocked_os):
        '''Test is_in_docker_container method for scenario
            when METAPARTICLE_IN_CONTAINER environment
            variable contains value '1'.
        '''
        expected_value = True
        mocked_os.getenv.return_value = '1'

        actual_value = containerize.is_in_docker_container()

        self.assertEqual(actual_value, expected_value)

    @patch("metaparticle_pkg.containerize.os")
    def test_is_in_docker_container_with_false_in_env(self, mocked_os):
        '''Test is_in_docker_container method for scenario
            when METAPARTICLE_IN_CONTAINER environment
            variable contains value 'false'.
        '''
        expected_value = False
        mocked_os.getenv.return_value = 'false'

        actual_value = containerize.is_in_docker_container()

        self.assertEqual(actual_value, expected_value)

    @patch("metaparticle_pkg.containerize.os")
    def test_is_in_docker_container_with_zero_in_env(self, mocked_os):
        '''Test is_in_docker_container method for scenario
            when METAPARTICLE_IN_CONTAINER environment
            variable contains value '0'.
        '''
        expected_value = False
        mocked_os.getenv.return_value = '0'

        actual_value = containerize.is_in_docker_container()

        self.assertEqual(actual_value, expected_value)

    @patch('metaparticle_pkg.containerize.os')
    def test_is_in_docker_container_with_docker_in_last_line(
            self, mocked_os):
        '''Test is_in_docker_container method
            in case of word "docker" in last line
        '''
        expected_value = True
        mocked_open_function = mock_open(read_data='line1\ndocker')
        mocked_os.getenv.return_value = None

        with patch("metaparticle_pkg.containerize.open", mocked_open_function):
            actual_value = containerize.is_in_docker_container()

            self.assertEqual(actual_value, expected_value)

    @patch('metaparticle_pkg.containerize.os')
    def test_is_in_docker_container_with_kubepods_in_last_line(
            self, mocked_os):
        '''Test is_in_docker_container method
            in case of word "kubepods" in last line
        '''
        expected_value = True
        mocked_open_function = mock_open(read_data='line1\nkubepods')
        mocked_os.getenv.return_value = None

        with patch("metaparticle_pkg.containerize.open", mocked_open_function):
            actual_value = containerize.is_in_docker_container()

            self.assertEqual(actual_value, expected_value)

    @patch('metaparticle_pkg.containerize.os')
    def test_is_in_docker_container_with_random_string_in_last_line(
            self, mocked_os):
        '''Test is_in_docker_container method in case of some randome
           string in last line
        '''
        expected_value = False
        mocked_open_function = mock_open(read_data='line1\nsome_random_string')
        mocked_os.getenv.return_value = None

        with patch("metaparticle_pkg.containerize.open", mocked_open_function):
            actual_value = containerize.is_in_docker_container()

            self.assertEqual(actual_value, expected_value)

    @patch('metaparticle_pkg.containerize.os')
    def test_is_in_docker_container_with_io_error(self, mocked_os):
        '''Test is_in_docker_container method in case of IOErrror'''
        expected_value = False
        mocked_open_function = mock_open()
        mocked_open_function.side_effect = IOError
        mocked_os.getenv.return_value = None

        with patch("metaparticle_pkg.containerize.open", mocked_open_function):
            actual_value = containerize.is_in_docker_container()
            self.assertEqual(actual_value, expected_value)

    @patch('metaparticle_pkg.containerize.shutil')
    def test_write_dockerfile(self, mocked_shutil):
        '''Test write_dockerfile method in case of Dockerfile is present'''

        # Input parameters
        package = MagicMock()
        package.dockerfile = '/some/fake_path_to_dockerfile'
        exec_file = '/some/fake_exec_file_path'

        containerize.write_dockerfile(package, exec_file)

        mocked_shutil.copy.assert_called_once_with(
            package.dockerfile, 'Dockerfile')

    def test_write_dockerfile_with_dockerfile_absent(self):
        '''Test write_dockerfile method in case of Dockerfile is absent'''

        mocked_open_function = mock_open()

        # Input parameters
        package = MagicMock()
        package.dockerfile = None
        exec_file = '/some/fake_exec_file_path'

        with patch("metaparticle_pkg.containerize.open",
                   mocked_open_function) as mocked_open:
            containerize.write_dockerfile(package, exec_file)
            self.assertEqual(mocked_open().write.call_count, 1)

    @patch("metaparticle_pkg.containerize.os")
    @patch("metaparticle_pkg.containerize.builder")
    @patch("metaparticle_pkg.containerize.runner")
    @patch("metaparticle_pkg.containerize.option")
    def test_containerize_decorator_one(
            self, mocked_option, mocked_runner, mocked_builder, mocked_os):
        '''Test Containerize decorator in case of
            is_in_docker_container return True
        '''
        mocked_os.getenv.return_value = 'true'
        expected_value = FunctionType

        # Input parameters
        def test_func(): pass

        actual_value = containerize.Containerize()(test_func)

        self.assertEqual(type(actual_value), expected_value)

    @patch("metaparticle_pkg.containerize.os")
    @patch("metaparticle_pkg.containerize.shutil")
    @patch("metaparticle_pkg.containerize.builder")
    @patch("metaparticle_pkg.containerize.runner")
    @patch("metaparticle_pkg.containerize.option")
    @patch("metaparticle_pkg.containerize.signal")
    def test_containerize_decorator_two(
        self, mocked_signal, mocked_option, mocked_runner,
            mocked_builder, mocked_shutil, mocked_os):
        '''Test Containerize decorator in case of
            is_in_docker_container return False
        '''
        mocked_os.getenv.return_value = 'false'

        # Input parameters
        def test_func(): pass

        containerize.Containerize()(test_func)()

        # Assert for calls in constructor
        self.assertEqual(mocked_option.load.call_count, 2)
        self.assertEqual(mocked_builder.select.call_count, 1)
        self.assertEqual(mocked_runner.select.call_count, 1)

        # Assert for calls in __call__ section
        self.assertEqual(mocked_os.getenv.call_count, 1)
        self.assertEqual(mocked_signal.signal.call_count, 1)
        self.assertEqual(mocked_runner.select().run.call_count, 1)
        self.assertEqual(mocked_runner.select().logs.call_count, 1)
        self.assertEqual(mocked_builder.select().build.call_count, 1)


if __name__ == '__main__':
    unittest.main()
