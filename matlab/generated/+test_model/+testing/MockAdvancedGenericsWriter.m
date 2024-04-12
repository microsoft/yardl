% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockAdvancedGenericsWriter < matlab.mixin.Copyable & test_model.AdvancedGenericsWriterBase
  properties
    testCase_
    expected_float_image_image
    expected_generic_record_1
    expected_tuple_of_optionals
    expected_tuple_of_optionals_alternate_syntax
    expected_tuple_of_vectors
  end

  methods
    function obj = MockAdvancedGenericsWriter(testCase)
      obj.testCase_ = testCase;
      obj.expected_float_image_image = yardl.None;
      obj.expected_generic_record_1 = yardl.None;
      obj.expected_tuple_of_optionals = yardl.None;
      obj.expected_tuple_of_optionals_alternate_syntax = yardl.None;
      obj.expected_tuple_of_vectors = yardl.None;
    end

    function expect_write_float_image_image_(obj, value)
      obj.expected_float_image_image = yardl.Optional(value);
    end

    function expect_write_generic_record_1_(obj, value)
      obj.expected_generic_record_1 = yardl.Optional(value);
    end

    function expect_write_tuple_of_optionals_(obj, value)
      obj.expected_tuple_of_optionals = yardl.Optional(value);
    end

    function expect_write_tuple_of_optionals_alternate_syntax_(obj, value)
      obj.expected_tuple_of_optionals_alternate_syntax = yardl.Optional(value);
    end

    function expect_write_tuple_of_vectors_(obj, value)
      obj.expected_tuple_of_vectors = yardl.Optional(value);
    end

    function verify(obj)
      obj.testCase_.verifyEqual(obj.expected_float_image_image, yardl.None, "Expected call to write_float_image_image_ was not received");
      obj.testCase_.verifyEqual(obj.expected_generic_record_1, yardl.None, "Expected call to write_generic_record_1_ was not received");
      obj.testCase_.verifyEqual(obj.expected_tuple_of_optionals, yardl.None, "Expected call to write_tuple_of_optionals_ was not received");
      obj.testCase_.verifyEqual(obj.expected_tuple_of_optionals_alternate_syntax, yardl.None, "Expected call to write_tuple_of_optionals_alternate_syntax_ was not received");
      obj.testCase_.verifyEqual(obj.expected_tuple_of_vectors, yardl.None, "Expected call to write_tuple_of_vectors_ was not received");
    end
  end

  methods (Access=protected)
    function write_float_image_image_(obj, value)
      obj.testCase_.verifyTrue(obj.expected_float_image_image.has_value(), "Unexpected call to write_float_image_image_");
      obj.testCase_.verifyEqual(value, obj.expected_float_image_image.value, "Unexpected argument value for call to write_float_image_image_");
      obj.expected_float_image_image = yardl.None;
    end

    function write_generic_record_1_(obj, value)
      obj.testCase_.verifyTrue(obj.expected_generic_record_1.has_value(), "Unexpected call to write_generic_record_1_");
      obj.testCase_.verifyEqual(value, obj.expected_generic_record_1.value, "Unexpected argument value for call to write_generic_record_1_");
      obj.expected_generic_record_1 = yardl.None;
    end

    function write_tuple_of_optionals_(obj, value)
      obj.testCase_.verifyTrue(obj.expected_tuple_of_optionals.has_value(), "Unexpected call to write_tuple_of_optionals_");
      obj.testCase_.verifyEqual(value, obj.expected_tuple_of_optionals.value, "Unexpected argument value for call to write_tuple_of_optionals_");
      obj.expected_tuple_of_optionals = yardl.None;
    end

    function write_tuple_of_optionals_alternate_syntax_(obj, value)
      obj.testCase_.verifyTrue(obj.expected_tuple_of_optionals_alternate_syntax.has_value(), "Unexpected call to write_tuple_of_optionals_alternate_syntax_");
      obj.testCase_.verifyEqual(value, obj.expected_tuple_of_optionals_alternate_syntax.value, "Unexpected argument value for call to write_tuple_of_optionals_alternate_syntax_");
      obj.expected_tuple_of_optionals_alternate_syntax = yardl.None;
    end

    function write_tuple_of_vectors_(obj, value)
      obj.testCase_.verifyTrue(obj.expected_tuple_of_vectors.has_value(), "Unexpected call to write_tuple_of_vectors_");
      obj.testCase_.verifyEqual(value, obj.expected_tuple_of_vectors.value, "Unexpected argument value for call to write_tuple_of_vectors_");
      obj.expected_tuple_of_vectors = yardl.None;
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
