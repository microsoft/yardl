% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockNDArraysWriter < matlab.mixin.Copyable & test_model.NDArraysWriterBase
  properties
    testCase_
    expected_ints
    expected_simple_record_array
    expected_record_with_vlens_array
    expected_record_with_nd_arrays
    expected_named_array
  end

  methods
    function self = MockNDArraysWriter(testCase)
      self.testCase_ = testCase;
      self.expected_ints = yardl.None;
      self.expected_simple_record_array = yardl.None;
      self.expected_record_with_vlens_array = yardl.None;
      self.expected_record_with_nd_arrays = yardl.None;
      self.expected_named_array = yardl.None;
    end

    function expect_write_ints_(self, value)
      self.expected_ints = yardl.Optional(value);
    end

    function expect_write_simple_record_array_(self, value)
      self.expected_simple_record_array = yardl.Optional(value);
    end

    function expect_write_record_with_vlens_array_(self, value)
      self.expected_record_with_vlens_array = yardl.Optional(value);
    end

    function expect_write_record_with_nd_arrays_(self, value)
      self.expected_record_with_nd_arrays = yardl.Optional(value);
    end

    function expect_write_named_array_(self, value)
      self.expected_named_array = yardl.Optional(value);
    end

    function verify(self)
      self.testCase_.verifyEqual(self.expected_ints, yardl.None, "Expected call to write_ints_ was not received");
      self.testCase_.verifyEqual(self.expected_simple_record_array, yardl.None, "Expected call to write_simple_record_array_ was not received");
      self.testCase_.verifyEqual(self.expected_record_with_vlens_array, yardl.None, "Expected call to write_record_with_vlens_array_ was not received");
      self.testCase_.verifyEqual(self.expected_record_with_nd_arrays, yardl.None, "Expected call to write_record_with_nd_arrays_ was not received");
      self.testCase_.verifyEqual(self.expected_named_array, yardl.None, "Expected call to write_named_array_ was not received");
    end
  end

  methods (Access=protected)
    function write_ints_(self, value)
      self.testCase_.verifyTrue(self.expected_ints.has_value(), "Unexpected call to write_ints_");
      self.testCase_.verifyEqual(value, self.expected_ints.value, "Unexpected argument value for call to write_ints_");
      self.expected_ints = yardl.None;
    end

    function write_simple_record_array_(self, value)
      self.testCase_.verifyTrue(self.expected_simple_record_array.has_value(), "Unexpected call to write_simple_record_array_");
      self.testCase_.verifyEqual(value, self.expected_simple_record_array.value, "Unexpected argument value for call to write_simple_record_array_");
      self.expected_simple_record_array = yardl.None;
    end

    function write_record_with_vlens_array_(self, value)
      self.testCase_.verifyTrue(self.expected_record_with_vlens_array.has_value(), "Unexpected call to write_record_with_vlens_array_");
      self.testCase_.verifyEqual(value, self.expected_record_with_vlens_array.value, "Unexpected argument value for call to write_record_with_vlens_array_");
      self.expected_record_with_vlens_array = yardl.None;
    end

    function write_record_with_nd_arrays_(self, value)
      self.testCase_.verifyTrue(self.expected_record_with_nd_arrays.has_value(), "Unexpected call to write_record_with_nd_arrays_");
      self.testCase_.verifyEqual(value, self.expected_record_with_nd_arrays.value, "Unexpected argument value for call to write_record_with_nd_arrays_");
      self.expected_record_with_nd_arrays = yardl.None;
    end

    function write_named_array_(self, value)
      self.testCase_.verifyTrue(self.expected_named_array.has_value(), "Unexpected call to write_named_array_");
      self.testCase_.verifyEqual(value, self.expected_named_array.value, "Unexpected argument value for call to write_named_array_");
      self.expected_named_array = yardl.None;
    end

    function close_(self)
    end
    function end_stream_(self)
    end
  end
end
