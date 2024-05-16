% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockSimpleGenericsWriter < matlab.mixin.Copyable & test_model.SimpleGenericsWriterBase
  properties
    testCase_
    expected_float_image
    expected_int_image
    expected_int_image_alternate_syntax
    expected_string_image
    expected_int_float_tuple
    expected_float_float_tuple
    expected_int_float_tuple_alternate_syntax
    expected_int_string_tuple
    expected_stream_of_type_variants
  end

  methods
    function self = MockSimpleGenericsWriter(testCase)
      self.testCase_ = testCase;
      self.expected_float_image = yardl.None;
      self.expected_int_image = yardl.None;
      self.expected_int_image_alternate_syntax = yardl.None;
      self.expected_string_image = yardl.None;
      self.expected_int_float_tuple = yardl.None;
      self.expected_float_float_tuple = yardl.None;
      self.expected_int_float_tuple_alternate_syntax = yardl.None;
      self.expected_int_string_tuple = yardl.None;
      self.expected_stream_of_type_variants = {};
    end

    function expect_write_float_image_(self, value)
      self.expected_float_image = yardl.Optional(value);
    end

    function expect_write_int_image_(self, value)
      self.expected_int_image = yardl.Optional(value);
    end

    function expect_write_int_image_alternate_syntax_(self, value)
      self.expected_int_image_alternate_syntax = yardl.Optional(value);
    end

    function expect_write_string_image_(self, value)
      self.expected_string_image = yardl.Optional(value);
    end

    function expect_write_int_float_tuple_(self, value)
      self.expected_int_float_tuple = yardl.Optional(value);
    end

    function expect_write_float_float_tuple_(self, value)
      self.expected_float_float_tuple = yardl.Optional(value);
    end

    function expect_write_int_float_tuple_alternate_syntax_(self, value)
      self.expected_int_float_tuple_alternate_syntax = yardl.Optional(value);
    end

    function expect_write_int_string_tuple_(self, value)
      self.expected_int_string_tuple = yardl.Optional(value);
    end

    function expect_write_stream_of_type_variants_(self, value)
      if iscell(value)
        for n = 1:numel(value)
          self.expected_stream_of_type_variants{end+1} = value{n};
        end
        return;
      end
      shape = size(value);
      lastDim = ndims(value);
      count = shape(lastDim);
      index = repelem({':'}, lastDim-1);
      for n = 1:count
        self.expected_stream_of_type_variants{end+1} = value(index{:}, n);
      end
    end

    function verify(self)
      self.testCase_.verifyEqual(self.expected_float_image, yardl.None, "Expected call to write_float_image_ was not received");
      self.testCase_.verifyEqual(self.expected_int_image, yardl.None, "Expected call to write_int_image_ was not received");
      self.testCase_.verifyEqual(self.expected_int_image_alternate_syntax, yardl.None, "Expected call to write_int_image_alternate_syntax_ was not received");
      self.testCase_.verifyEqual(self.expected_string_image, yardl.None, "Expected call to write_string_image_ was not received");
      self.testCase_.verifyEqual(self.expected_int_float_tuple, yardl.None, "Expected call to write_int_float_tuple_ was not received");
      self.testCase_.verifyEqual(self.expected_float_float_tuple, yardl.None, "Expected call to write_float_float_tuple_ was not received");
      self.testCase_.verifyEqual(self.expected_int_float_tuple_alternate_syntax, yardl.None, "Expected call to write_int_float_tuple_alternate_syntax_ was not received");
      self.testCase_.verifyEqual(self.expected_int_string_tuple, yardl.None, "Expected call to write_int_string_tuple_ was not received");
      self.testCase_.verifyTrue(isempty(self.expected_stream_of_type_variants), "Expected call to write_stream_of_type_variants_ was not received");
    end
  end

  methods (Access=protected)
    function write_float_image_(self, value)
      self.testCase_.verifyTrue(self.expected_float_image.has_value(), "Unexpected call to write_float_image_");
      self.testCase_.verifyEqual(value, self.expected_float_image.value, "Unexpected argument value for call to write_float_image_");
      self.expected_float_image = yardl.None;
    end

    function write_int_image_(self, value)
      self.testCase_.verifyTrue(self.expected_int_image.has_value(), "Unexpected call to write_int_image_");
      self.testCase_.verifyEqual(value, self.expected_int_image.value, "Unexpected argument value for call to write_int_image_");
      self.expected_int_image = yardl.None;
    end

    function write_int_image_alternate_syntax_(self, value)
      self.testCase_.verifyTrue(self.expected_int_image_alternate_syntax.has_value(), "Unexpected call to write_int_image_alternate_syntax_");
      self.testCase_.verifyEqual(value, self.expected_int_image_alternate_syntax.value, "Unexpected argument value for call to write_int_image_alternate_syntax_");
      self.expected_int_image_alternate_syntax = yardl.None;
    end

    function write_string_image_(self, value)
      self.testCase_.verifyTrue(self.expected_string_image.has_value(), "Unexpected call to write_string_image_");
      self.testCase_.verifyEqual(value, self.expected_string_image.value, "Unexpected argument value for call to write_string_image_");
      self.expected_string_image = yardl.None;
    end

    function write_int_float_tuple_(self, value)
      self.testCase_.verifyTrue(self.expected_int_float_tuple.has_value(), "Unexpected call to write_int_float_tuple_");
      self.testCase_.verifyEqual(value, self.expected_int_float_tuple.value, "Unexpected argument value for call to write_int_float_tuple_");
      self.expected_int_float_tuple = yardl.None;
    end

    function write_float_float_tuple_(self, value)
      self.testCase_.verifyTrue(self.expected_float_float_tuple.has_value(), "Unexpected call to write_float_float_tuple_");
      self.testCase_.verifyEqual(value, self.expected_float_float_tuple.value, "Unexpected argument value for call to write_float_float_tuple_");
      self.expected_float_float_tuple = yardl.None;
    end

    function write_int_float_tuple_alternate_syntax_(self, value)
      self.testCase_.verifyTrue(self.expected_int_float_tuple_alternate_syntax.has_value(), "Unexpected call to write_int_float_tuple_alternate_syntax_");
      self.testCase_.verifyEqual(value, self.expected_int_float_tuple_alternate_syntax.value, "Unexpected argument value for call to write_int_float_tuple_alternate_syntax_");
      self.expected_int_float_tuple_alternate_syntax = yardl.None;
    end

    function write_int_string_tuple_(self, value)
      self.testCase_.verifyTrue(self.expected_int_string_tuple.has_value(), "Unexpected call to write_int_string_tuple_");
      self.testCase_.verifyEqual(value, self.expected_int_string_tuple.value, "Unexpected argument value for call to write_int_string_tuple_");
      self.expected_int_string_tuple = yardl.None;
    end

    function write_stream_of_type_variants_(self, value)
      assert(iscell(value));
      assert(isscalar(value));
      self.testCase_.verifyFalse(isempty(self.expected_stream_of_type_variants), "Unexpected call to write_stream_of_type_variants_");
      self.testCase_.verifyEqual(value{1}, self.expected_stream_of_type_variants{1}, "Unexpected argument value for call to write_stream_of_type_variants_");
      self.expected_stream_of_type_variants = self.expected_stream_of_type_variants(2:end);
    end

    function close_(self)
    end
    function end_stream_(self)
    end
  end
end
