% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockAliasesWriter < matlab.mixin.Copyable & test_model.AliasesWriterBase
  properties
    testCase_
    expected_aliased_string
    expected_aliased_enum
    expected_aliased_open_generic
    expected_aliased_closed_generic
    expected_aliased_optional
    expected_aliased_generic_optional
    expected_aliased_generic_union_2
    expected_aliased_generic_vector
    expected_aliased_generic_fixed_vector
    expected_stream_of_aliased_generic_union_2
  end

  methods
    function self = MockAliasesWriter(testCase)
      self.testCase_ = testCase;
      self.expected_aliased_string = yardl.None;
      self.expected_aliased_enum = yardl.None;
      self.expected_aliased_open_generic = yardl.None;
      self.expected_aliased_closed_generic = yardl.None;
      self.expected_aliased_optional = yardl.None;
      self.expected_aliased_generic_optional = yardl.None;
      self.expected_aliased_generic_union_2 = yardl.None;
      self.expected_aliased_generic_vector = yardl.None;
      self.expected_aliased_generic_fixed_vector = yardl.None;
      self.expected_stream_of_aliased_generic_union_2 = {};
    end

    function expect_write_aliased_string_(self, value)
      self.expected_aliased_string = yardl.Optional(value);
    end

    function expect_write_aliased_enum_(self, value)
      self.expected_aliased_enum = yardl.Optional(value);
    end

    function expect_write_aliased_open_generic_(self, value)
      self.expected_aliased_open_generic = yardl.Optional(value);
    end

    function expect_write_aliased_closed_generic_(self, value)
      self.expected_aliased_closed_generic = yardl.Optional(value);
    end

    function expect_write_aliased_optional_(self, value)
      self.expected_aliased_optional = yardl.Optional(value);
    end

    function expect_write_aliased_generic_optional_(self, value)
      self.expected_aliased_generic_optional = yardl.Optional(value);
    end

    function expect_write_aliased_generic_union_2_(self, value)
      self.expected_aliased_generic_union_2 = yardl.Optional(value);
    end

    function expect_write_aliased_generic_vector_(self, value)
      self.expected_aliased_generic_vector = yardl.Optional(value);
    end

    function expect_write_aliased_generic_fixed_vector_(self, value)
      self.expected_aliased_generic_fixed_vector = yardl.Optional(value);
    end

    function expect_write_stream_of_aliased_generic_union_2_(self, value)
      if iscell(value)
        for n = 1:numel(value)
          self.expected_stream_of_aliased_generic_union_2{end+1} = value{n};
        end
        return;
      end
      shape = size(value);
      lastDim = ndims(value);
      count = shape(lastDim);
      index = repelem({':'}, lastDim-1);
      for n = 1:count
        self.expected_stream_of_aliased_generic_union_2{end+1} = value(index{:}, n);
      end
    end

    function verify(self)
      self.testCase_.verifyEqual(self.expected_aliased_string, yardl.None, "Expected call to write_aliased_string_ was not received");
      self.testCase_.verifyEqual(self.expected_aliased_enum, yardl.None, "Expected call to write_aliased_enum_ was not received");
      self.testCase_.verifyEqual(self.expected_aliased_open_generic, yardl.None, "Expected call to write_aliased_open_generic_ was not received");
      self.testCase_.verifyEqual(self.expected_aliased_closed_generic, yardl.None, "Expected call to write_aliased_closed_generic_ was not received");
      self.testCase_.verifyEqual(self.expected_aliased_optional, yardl.None, "Expected call to write_aliased_optional_ was not received");
      self.testCase_.verifyEqual(self.expected_aliased_generic_optional, yardl.None, "Expected call to write_aliased_generic_optional_ was not received");
      self.testCase_.verifyEqual(self.expected_aliased_generic_union_2, yardl.None, "Expected call to write_aliased_generic_union_2_ was not received");
      self.testCase_.verifyEqual(self.expected_aliased_generic_vector, yardl.None, "Expected call to write_aliased_generic_vector_ was not received");
      self.testCase_.verifyEqual(self.expected_aliased_generic_fixed_vector, yardl.None, "Expected call to write_aliased_generic_fixed_vector_ was not received");
      self.testCase_.verifyTrue(isempty(self.expected_stream_of_aliased_generic_union_2), "Expected call to write_stream_of_aliased_generic_union_2_ was not received");
    end
  end

  methods (Access=protected)
    function write_aliased_string_(self, value)
      self.testCase_.verifyTrue(self.expected_aliased_string.has_value(), "Unexpected call to write_aliased_string_");
      self.testCase_.verifyEqual(value, self.expected_aliased_string.value, "Unexpected argument value for call to write_aliased_string_");
      self.expected_aliased_string = yardl.None;
    end

    function write_aliased_enum_(self, value)
      self.testCase_.verifyTrue(self.expected_aliased_enum.has_value(), "Unexpected call to write_aliased_enum_");
      self.testCase_.verifyEqual(value, self.expected_aliased_enum.value, "Unexpected argument value for call to write_aliased_enum_");
      self.expected_aliased_enum = yardl.None;
    end

    function write_aliased_open_generic_(self, value)
      self.testCase_.verifyTrue(self.expected_aliased_open_generic.has_value(), "Unexpected call to write_aliased_open_generic_");
      self.testCase_.verifyEqual(value, self.expected_aliased_open_generic.value, "Unexpected argument value for call to write_aliased_open_generic_");
      self.expected_aliased_open_generic = yardl.None;
    end

    function write_aliased_closed_generic_(self, value)
      self.testCase_.verifyTrue(self.expected_aliased_closed_generic.has_value(), "Unexpected call to write_aliased_closed_generic_");
      self.testCase_.verifyEqual(value, self.expected_aliased_closed_generic.value, "Unexpected argument value for call to write_aliased_closed_generic_");
      self.expected_aliased_closed_generic = yardl.None;
    end

    function write_aliased_optional_(self, value)
      self.testCase_.verifyTrue(self.expected_aliased_optional.has_value(), "Unexpected call to write_aliased_optional_");
      self.testCase_.verifyEqual(value, self.expected_aliased_optional.value, "Unexpected argument value for call to write_aliased_optional_");
      self.expected_aliased_optional = yardl.None;
    end

    function write_aliased_generic_optional_(self, value)
      self.testCase_.verifyTrue(self.expected_aliased_generic_optional.has_value(), "Unexpected call to write_aliased_generic_optional_");
      self.testCase_.verifyEqual(value, self.expected_aliased_generic_optional.value, "Unexpected argument value for call to write_aliased_generic_optional_");
      self.expected_aliased_generic_optional = yardl.None;
    end

    function write_aliased_generic_union_2_(self, value)
      self.testCase_.verifyTrue(self.expected_aliased_generic_union_2.has_value(), "Unexpected call to write_aliased_generic_union_2_");
      self.testCase_.verifyEqual(value, self.expected_aliased_generic_union_2.value, "Unexpected argument value for call to write_aliased_generic_union_2_");
      self.expected_aliased_generic_union_2 = yardl.None;
    end

    function write_aliased_generic_vector_(self, value)
      self.testCase_.verifyTrue(self.expected_aliased_generic_vector.has_value(), "Unexpected call to write_aliased_generic_vector_");
      self.testCase_.verifyEqual(value, self.expected_aliased_generic_vector.value, "Unexpected argument value for call to write_aliased_generic_vector_");
      self.expected_aliased_generic_vector = yardl.None;
    end

    function write_aliased_generic_fixed_vector_(self, value)
      self.testCase_.verifyTrue(self.expected_aliased_generic_fixed_vector.has_value(), "Unexpected call to write_aliased_generic_fixed_vector_");
      self.testCase_.verifyEqual(value, self.expected_aliased_generic_fixed_vector.value, "Unexpected argument value for call to write_aliased_generic_fixed_vector_");
      self.expected_aliased_generic_fixed_vector = yardl.None;
    end

    function write_stream_of_aliased_generic_union_2_(self, value)
      assert(iscell(value));
      assert(isscalar(value));
      self.testCase_.verifyFalse(isempty(self.expected_stream_of_aliased_generic_union_2), "Unexpected call to write_stream_of_aliased_generic_union_2_");
      self.testCase_.verifyEqual(value{1}, self.expected_stream_of_aliased_generic_union_2{1}, "Unexpected argument value for call to write_stream_of_aliased_generic_union_2_");
      self.expected_stream_of_aliased_generic_union_2 = self.expected_stream_of_aliased_generic_union_2(2:end);
    end

    function close_(self)
    end
    function end_stream_(self)
    end
  end
end
