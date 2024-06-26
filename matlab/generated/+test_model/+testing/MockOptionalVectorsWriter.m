% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockOptionalVectorsWriter < matlab.mixin.Copyable & test_model.OptionalVectorsWriterBase
  properties
    testCase_
    expected_record_with_optional_vector
  end

  methods
    function self = MockOptionalVectorsWriter(testCase)
      self.testCase_ = testCase;
      self.expected_record_with_optional_vector = yardl.None;
    end

    function expect_write_record_with_optional_vector_(self, value)
      self.expected_record_with_optional_vector = yardl.Optional(value);
    end

    function verify(self)
      self.testCase_.verifyEqual(self.expected_record_with_optional_vector, yardl.None, "Expected call to write_record_with_optional_vector_ was not received");
    end
  end

  methods (Access=protected)
    function write_record_with_optional_vector_(self, value)
      self.testCase_.verifyTrue(self.expected_record_with_optional_vector.has_value(), "Unexpected call to write_record_with_optional_vector_");
      self.testCase_.verifyEqual(value, self.expected_record_with_optional_vector.value, "Unexpected argument value for call to write_record_with_optional_vector_");
      self.expected_record_with_optional_vector = yardl.None;
    end

    function close_(self)
    end
    function end_stream_(self)
    end
  end
end
