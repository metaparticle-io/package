#!/usr/bin/env python
'''Unit tests for MetaparticleRunner'''

import unittest
from unittest.mock import patch, MagicMock, mock_open
from metaparticle_pkg.runner import metaparticle


class TestMetaparticleRunner(unittest.TestCase):
    '''Unit tests for MetaparticleRunner'''

    def setUp(self):
        self.metaparticle_runner = metaparticle.MetaparticleRunner()

    @patch("metaparticle_pkg.runner.metaparticle.subprocess")
    def test_cancel(self, mocked_subprocess):
        '''Test cancel method'''

        # Input arguments
        name = "test_name"

        # Expected argument called with
        expected_args = [
            'mp-compiler', '-f', '.metaparticle/spec.json', '--delete']

        self.metaparticle_runner.cancel(name)

        mocked_subprocess.check_call.assert_called_once_with(
            expected_args)

    @patch("metaparticle_pkg.runner.metaparticle.subprocess")
    def test_logs(self, mocked_subprocess):
        '''Test logs method'''

        # Input arguments
        name = "test_name"

        # Expected argument called with
        expected_args = [
            'mp-compiler', '-f', '.metaparticle/spec.json',
            '--deploy=false', '--attach=true'
        ]

        self.metaparticle_runner.logs(name)

        mocked_subprocess.check_call.assert_called_once_with(
            expected_args)

    def test_ports(self):
        '''Test ports method'''

        # Input arguments
        ports = ['4686']

        # Expected result
        expected_result = [
            {'number': '4686', 'protocol': 'TCP'}
        ]

        actual_result = self.metaparticle_runner.ports(ports)

        self.assertEqual(expected_result, actual_result)

    @patch("metaparticle_pkg.runner.metaparticle.subprocess")
    @patch("metaparticle_pkg.runner.metaparticle.os")
    @patch("metaparticle_pkg.runner.metaparticle.json")
    def test_run_when_dir_not_exists(
            self, mocked_json, mocked_os, mocked_subprocess):
        '''Test run method in scenario when directory
           .metaparticle does not exist
        '''
        mocked_os.path.exists.return_value = False
        mocked_open_function = mock_open()

        # Input parameters
        img_name = "test_image"
        name = "test_name"
        options = MagicMock()
        options.ports = ['23235']
        options.shardSpec = 1
        options.replicas = 1

        # Mock built-in open function
        with patch("metaparticle_pkg.runner.metaparticle.open",
                   mocked_open_function):
            self.metaparticle_runner.run(img_name, name, options)

            mocked_os.path.exists.assert_called_once_with(
                '.metaparticle')
            mocked_os.makedirs.assert_called_once_with(
                '.metaparticle')
            self.assertEqual(mocked_json.dump.call_count, 1)
            mocked_subprocess.check_call.assert_called_once_with(
                ['mp-compiler', '-f', '.metaparticle/spec.json']
            )

    @patch("metaparticle_pkg.runner.metaparticle.subprocess")
    @patch("metaparticle_pkg.runner.metaparticle.os")
    @patch("metaparticle_pkg.runner.metaparticle.json")
    def test_run_when_dir_exists(
            self, mocked_json, mocked_os, mocked_subprocess):
        '''Test run method in scenario when directory
           .metaparticle exists
        '''
        mocked_os.path.exists.return_value = True
        mocked_open_function = mock_open()

        # Input parameters
        img_name = "test_image"
        name = "test_name"
        options = MagicMock()
        options.ports = ['23235']
        options.shardSpec = 1
        options.replicas = 1

        with patch("metaparticle_pkg.runner.metaparticle.open",
                   mocked_open_function):
            self.metaparticle_runner.run(img_name, name, options)

            mocked_os.path.exists.assert_called_once_with(
                '.metaparticle')
            self.assertFalse(mocked_os.makedirs.called)
            self.assertEqual(mocked_json.dump.call_count, 1)
            mocked_subprocess.check_call.assert_called_once_with(
                ['mp-compiler', '-f', '.metaparticle/spec.json']
            )


if __name__ == '__main__':
    unittest.main()
