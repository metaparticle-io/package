#!/usr/bin/env python
'''Unit tests for option module'''

from sys import exit as real_exit
import unittest
from unittest.mock import patch, MagicMock
from metaparticle_pkg import option


class TestOption(unittest.TestCase):
    '''Unit tests for option module'''

    @patch('metaparticle_pkg.option.sys')
    def test_load_for_non_dict_type(self, mocked_sys):
        '''Test load method in scenario when
           options parameter is not a dict type
        '''
        def fake_exit(exit_code):
            '''Function to mimick sys.exit
                with extra functionality
            '''
            self.assertTrue(mocked_sys.stderr.write.call_count, 1)
            real_exit(exit_code)

        mocked_sys.exit.side_effect = fake_exit

        # Input parameters
        test_option = []
        test_cls = None

        with self.assertRaises(SystemExit):
            option.load(test_cls, test_option)

    @patch('metaparticle_pkg.option.sys')
    def test_load_for_dict_type_with_missing_req_field(self, mocked_sys):
        '''Test load method in scenario when
           options parameter is a dict type
           and required_options contains option
           which is not part of options dictionary.
           It is a case of missing required field.
        '''
        def fake_exit(exit_code):
            '''Function to mimick sys.exit
                with extra functionality
            '''
            self.assertTrue(mocked_sys.stderr.write.call_count, 1)
            real_exit(exit_code)

        mocked_sys.exit.side_effect = fake_exit

        # Input parameters
        test_options = {}
        test_cls = MagicMock()
        test_cls.required_options = ["test_option"]

        with self.assertRaises(SystemExit):
            option.load(test_cls, test_options)

    def test_load_dict_type_with_req_field_available(self):
        '''Test load method in scenario when
           options parameter is a dict type
           and required_options contains option
           which is a part of options dictionary.
        '''
        # Input parameters
        test_options = {"test_option": "test_value"}
        test_cls = MagicMock
        test_cls.required_options = ["test_option"]

        actual_result = option.load(test_cls, test_options)
        self.assertEqual(
            actual_result.required_options, test_cls.required_options)

    @patch('metaparticle_pkg.option.sys')
    def test_load_for_type_error(self, mocked_sys):
        '''Test load method in case of TypeError.
        '''
        def fake_exit(exit_code):
            '''Function to mimick as sys.exit
                with extra functionality
            '''
            self.assertTrue(mocked_sys.stderr.write.call_count, 1)
            real_exit(exit_code)

        mocked_sys.exit.side_effect = fake_exit

        # Input parameters
        test_options = {}
        test_cls = MagicMock
        test_cls.side_effect = TypeError
        test_cls.required_options = ["test_option"]

        with self.assertRaises(TypeError):
            option.load(test_cls, test_options)


if __name__ == '__main__':
    unittest.main()
