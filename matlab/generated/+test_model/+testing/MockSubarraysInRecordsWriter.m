% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockSubarraysInRecordsWriter < matlab.mixin.Copyable & test_model.SubarraysInRecordsWriterBase
  properties
    testCase_
    expected_with_fixed_subarrays
    expected_with_vlen_subarrays
  end

  methods
    function self = MockSubarraysInRecordsWriter(testCase)
      self.testCase_ = testCase;
      self.expected_with_fixed_subarrays = yardl.None;
      self.expected_with_vlen_subarrays = yardl.None;
    end

    function expect_write_with_fixed_subarrays_(self, value)
      self.expected_with_fixed_subarrays = yardl.Optional(value);
    end

    function expect_write_with_vlen_subarrays_(self, value)
      self.expected_with_vlen_subarrays = yardl.Optional(value);
    end

    function verify(self)
      self.testCase_.verifyEqual(self.expected_with_fixed_subarrays, yardl.None, "Expected call to write_with_fixed_subarrays_ was not received");
      self.testCase_.verifyEqual(self.expected_with_vlen_subarrays, yardl.None, "Expected call to write_with_vlen_subarrays_ was not received");
    end
  end

  methods (Access=protected)
    function write_with_fixed_subarrays_(self, value)
      self.testCase_.verifyTrue(self.expected_with_fixed_subarrays.has_value(), "Unexpected call to write_with_fixed_subarrays_");
      self.testCase_.verifyEqual(value, self.expected_with_fixed_subarrays.value, "Unexpected argument value for call to write_with_fixed_subarrays_");
      self.expected_with_fixed_subarrays = yardl.None;
    end

    function write_with_vlen_subarrays_(self, value)
      self.testCase_.verifyTrue(self.expected_with_vlen_subarrays.has_value(), "Unexpected call to write_with_vlen_subarrays_");
      self.testCase_.verifyEqual(value, self.expected_with_vlen_subarrays.value, "Unexpected argument value for call to write_with_vlen_subarrays_");
      self.expected_with_vlen_subarrays = yardl.None;
    end

    function close_(self)
    end
    function end_stream_(self)
    end
  end
end
