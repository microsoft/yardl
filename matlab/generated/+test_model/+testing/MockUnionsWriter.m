% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockUnionsWriter < matlab.mixin.Copyable & test_model.UnionsWriterBase
  properties
    testCase_
    expected_int_or_simple_record
    expected_int_or_record_with_vlens
    expected_monosotate_or_int_or_simple_record
    expected_record_with_unions
  end

  methods
    function self = MockUnionsWriter(testCase)
      self.testCase_ = testCase;
      self.expected_int_or_simple_record = yardl.None;
      self.expected_int_or_record_with_vlens = yardl.None;
      self.expected_monosotate_or_int_or_simple_record = yardl.None;
      self.expected_record_with_unions = yardl.None;
    end

    function expect_write_int_or_simple_record_(self, value)
      self.expected_int_or_simple_record = yardl.Optional(value);
    end

    function expect_write_int_or_record_with_vlens_(self, value)
      self.expected_int_or_record_with_vlens = yardl.Optional(value);
    end

    function expect_write_monosotate_or_int_or_simple_record_(self, value)
      self.expected_monosotate_or_int_or_simple_record = yardl.Optional(value);
    end

    function expect_write_record_with_unions_(self, value)
      self.expected_record_with_unions = yardl.Optional(value);
    end

    function verify(self)
      self.testCase_.verifyEqual(self.expected_int_or_simple_record, yardl.None, "Expected call to write_int_or_simple_record_ was not received");
      self.testCase_.verifyEqual(self.expected_int_or_record_with_vlens, yardl.None, "Expected call to write_int_or_record_with_vlens_ was not received");
      self.testCase_.verifyEqual(self.expected_monosotate_or_int_or_simple_record, yardl.None, "Expected call to write_monosotate_or_int_or_simple_record_ was not received");
      self.testCase_.verifyEqual(self.expected_record_with_unions, yardl.None, "Expected call to write_record_with_unions_ was not received");
    end
  end

  methods (Access=protected)
    function write_int_or_simple_record_(self, value)
      self.testCase_.verifyTrue(self.expected_int_or_simple_record.has_value(), "Unexpected call to write_int_or_simple_record_");
      self.testCase_.verifyEqual(value, self.expected_int_or_simple_record.value, "Unexpected argument value for call to write_int_or_simple_record_");
      self.expected_int_or_simple_record = yardl.None;
    end

    function write_int_or_record_with_vlens_(self, value)
      self.testCase_.verifyTrue(self.expected_int_or_record_with_vlens.has_value(), "Unexpected call to write_int_or_record_with_vlens_");
      self.testCase_.verifyEqual(value, self.expected_int_or_record_with_vlens.value, "Unexpected argument value for call to write_int_or_record_with_vlens_");
      self.expected_int_or_record_with_vlens = yardl.None;
    end

    function write_monosotate_or_int_or_simple_record_(self, value)
      self.testCase_.verifyTrue(self.expected_monosotate_or_int_or_simple_record.has_value(), "Unexpected call to write_monosotate_or_int_or_simple_record_");
      self.testCase_.verifyEqual(value, self.expected_monosotate_or_int_or_simple_record.value, "Unexpected argument value for call to write_monosotate_or_int_or_simple_record_");
      self.expected_monosotate_or_int_or_simple_record = yardl.None;
    end

    function write_record_with_unions_(self, value)
      self.testCase_.verifyTrue(self.expected_record_with_unions.has_value(), "Unexpected call to write_record_with_unions_");
      self.testCase_.verifyEqual(value, self.expected_record_with_unions.value, "Unexpected argument value for call to write_record_with_unions_");
      self.expected_record_with_unions = yardl.None;
    end

    function close_(self)
    end
    function end_stream_(self)
    end
  end
end
