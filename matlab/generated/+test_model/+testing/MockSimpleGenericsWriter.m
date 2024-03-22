classdef MockSimpleGenericsWriter < test_model.SimpleGenericsWriterBase
  properties
    testCase_
    write_float_image_written
    write_int_image_written
    write_int_image_alternate_syntax_written
    write_string_image_written
    write_int_float_tuple_written
    write_float_float_tuple_written
    write_int_float_tuple_alternate_syntax_written
    write_int_string_tuple_written
    write_stream_of_type_variants_written
  end

  methods
    function obj = MockSimpleGenericsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_float_image_written = Node.empty();
      obj.write_int_image_written = Node.empty();
      obj.write_int_image_alternate_syntax_written = Node.empty();
      obj.write_string_image_written = Node.empty();
      obj.write_int_float_tuple_written = Node.empty();
      obj.write_float_float_tuple_written = Node.empty();
      obj.write_int_float_tuple_alternate_syntax_written = Node.empty();
      obj.write_int_string_tuple_written = Node.empty();
      obj.write_stream_of_type_variants_written = Node.empty();
    end

    function expect_write_float_image_(obj, value)
      obj.write_float_image_written(end+1) = Node(value);
    end

    function expect_write_int_image_(obj, value)
      obj.write_int_image_written(end+1) = Node(value);
    end

    function expect_write_int_image_alternate_syntax_(obj, value)
      obj.write_int_image_alternate_syntax_written(end+1) = Node(value);
    end

    function expect_write_string_image_(obj, value)
      obj.write_string_image_written(end+1) = Node(value);
    end

    function expect_write_int_float_tuple_(obj, value)
      obj.write_int_float_tuple_written(end+1) = Node(value);
    end

    function expect_write_float_float_tuple_(obj, value)
      obj.write_float_float_tuple_written(end+1) = Node(value);
    end

    function expect_write_int_float_tuple_alternate_syntax_(obj, value)
      obj.write_int_float_tuple_alternate_syntax_written(end+1) = Node(value);
    end

    function expect_write_int_string_tuple_(obj, value)
      obj.write_int_string_tuple_written(end+1) = Node(value);
    end

    function expect_write_stream_of_type_variants_(obj, value)
      obj.write_stream_of_type_variants_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_float_image_written), "Expected call to write_float_image_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_int_image_written), "Expected call to write_int_image_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_int_image_alternate_syntax_written), "Expected call to write_int_image_alternate_syntax_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_string_image_written), "Expected call to write_string_image_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_int_float_tuple_written), "Expected call to write_int_float_tuple_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_float_float_tuple_written), "Expected call to write_float_float_tuple_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_int_float_tuple_alternate_syntax_written), "Expected call to write_int_float_tuple_alternate_syntax_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_int_string_tuple_written), "Expected call to write_int_string_tuple_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_stream_of_type_variants_written), "Expected call to write_stream_of_type_variants_ was not received");
    end
  end

  methods (Access=protected)
    function write_float_image_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_float_image_written), "Unexpected call to write_float_image_");
      expected = obj.write_float_image_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_float_image_");
      obj.write_float_image_written = obj.write_float_image_written(2:end);
    end

    function write_int_image_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_int_image_written), "Unexpected call to write_int_image_");
      expected = obj.write_int_image_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_int_image_");
      obj.write_int_image_written = obj.write_int_image_written(2:end);
    end

    function write_int_image_alternate_syntax_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_int_image_alternate_syntax_written), "Unexpected call to write_int_image_alternate_syntax_");
      expected = obj.write_int_image_alternate_syntax_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_int_image_alternate_syntax_");
      obj.write_int_image_alternate_syntax_written = obj.write_int_image_alternate_syntax_written(2:end);
    end

    function write_string_image_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_string_image_written), "Unexpected call to write_string_image_");
      expected = obj.write_string_image_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_string_image_");
      obj.write_string_image_written = obj.write_string_image_written(2:end);
    end

    function write_int_float_tuple_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_int_float_tuple_written), "Unexpected call to write_int_float_tuple_");
      expected = obj.write_int_float_tuple_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_int_float_tuple_");
      obj.write_int_float_tuple_written = obj.write_int_float_tuple_written(2:end);
    end

    function write_float_float_tuple_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_float_float_tuple_written), "Unexpected call to write_float_float_tuple_");
      expected = obj.write_float_float_tuple_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_float_float_tuple_");
      obj.write_float_float_tuple_written = obj.write_float_float_tuple_written(2:end);
    end

    function write_int_float_tuple_alternate_syntax_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_int_float_tuple_alternate_syntax_written), "Unexpected call to write_int_float_tuple_alternate_syntax_");
      expected = obj.write_int_float_tuple_alternate_syntax_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_int_float_tuple_alternate_syntax_");
      obj.write_int_float_tuple_alternate_syntax_written = obj.write_int_float_tuple_alternate_syntax_written(2:end);
    end

    function write_int_string_tuple_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_int_string_tuple_written), "Unexpected call to write_int_string_tuple_");
      expected = obj.write_int_string_tuple_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_int_string_tuple_");
      obj.write_int_string_tuple_written = obj.write_int_string_tuple_written(2:end);
    end

    function write_stream_of_type_variants_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_stream_of_type_variants_written), "Unexpected call to write_stream_of_type_variants_");
      expected = obj.write_stream_of_type_variants_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_stream_of_type_variants_");
      obj.write_stream_of_type_variants_written = obj.write_stream_of_type_variants_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end