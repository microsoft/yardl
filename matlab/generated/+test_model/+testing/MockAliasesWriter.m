classdef MockAliasesWriter < test_model.AliasesWriterBase
  properties
    testCase_
    write_aliased_string_written
    write_aliased_enum_written
    write_aliased_open_generic_written
    write_aliased_closed_generic_written
    write_aliased_optional_written
    write_aliased_generic_optional_written
    write_aliased_generic_union_2_written
    write_aliased_generic_vector_written
    write_aliased_generic_fixed_vector_written
    write_stream_of_aliased_generic_union_2_written
  end

  methods
    function obj = MockAliasesWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_aliased_string_written = Node.empty();
      obj.write_aliased_enum_written = Node.empty();
      obj.write_aliased_open_generic_written = Node.empty();
      obj.write_aliased_closed_generic_written = Node.empty();
      obj.write_aliased_optional_written = Node.empty();
      obj.write_aliased_generic_optional_written = Node.empty();
      obj.write_aliased_generic_union_2_written = Node.empty();
      obj.write_aliased_generic_vector_written = Node.empty();
      obj.write_aliased_generic_fixed_vector_written = Node.empty();
      obj.write_stream_of_aliased_generic_union_2_written = Node.empty();
    end

    function expect_write_aliased_string_(obj, value)
      obj.write_aliased_string_written(end+1) = Node(value);
    end

    function expect_write_aliased_enum_(obj, value)
      obj.write_aliased_enum_written(end+1) = Node(value);
    end

    function expect_write_aliased_open_generic_(obj, value)
      obj.write_aliased_open_generic_written(end+1) = Node(value);
    end

    function expect_write_aliased_closed_generic_(obj, value)
      obj.write_aliased_closed_generic_written(end+1) = Node(value);
    end

    function expect_write_aliased_optional_(obj, value)
      obj.write_aliased_optional_written(end+1) = Node(value);
    end

    function expect_write_aliased_generic_optional_(obj, value)
      obj.write_aliased_generic_optional_written(end+1) = Node(value);
    end

    function expect_write_aliased_generic_union_2_(obj, value)
      obj.write_aliased_generic_union_2_written(end+1) = Node(value);
    end

    function expect_write_aliased_generic_vector_(obj, value)
      obj.write_aliased_generic_vector_written(end+1) = Node(value);
    end

    function expect_write_aliased_generic_fixed_vector_(obj, value)
      obj.write_aliased_generic_fixed_vector_written(end+1) = Node(value);
    end

    function expect_write_stream_of_aliased_generic_union_2_(obj, value)
      obj.write_stream_of_aliased_generic_union_2_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_aliased_string_written), "Expected call to write_aliased_string_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_aliased_enum_written), "Expected call to write_aliased_enum_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_aliased_open_generic_written), "Expected call to write_aliased_open_generic_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_aliased_closed_generic_written), "Expected call to write_aliased_closed_generic_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_aliased_optional_written), "Expected call to write_aliased_optional_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_aliased_generic_optional_written), "Expected call to write_aliased_generic_optional_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_aliased_generic_union_2_written), "Expected call to write_aliased_generic_union_2_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_aliased_generic_vector_written), "Expected call to write_aliased_generic_vector_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_aliased_generic_fixed_vector_written), "Expected call to write_aliased_generic_fixed_vector_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_stream_of_aliased_generic_union_2_written), "Expected call to write_stream_of_aliased_generic_union_2_ was not received");
    end
  end

  methods (Access=protected)
    function write_aliased_string_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_aliased_string_written), "Unexpected call to write_aliased_string_");
      expected = obj.write_aliased_string_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_aliased_string_");
      obj.write_aliased_string_written = obj.write_aliased_string_written(2:end);
    end

    function write_aliased_enum_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_aliased_enum_written), "Unexpected call to write_aliased_enum_");
      expected = obj.write_aliased_enum_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_aliased_enum_");
      obj.write_aliased_enum_written = obj.write_aliased_enum_written(2:end);
    end

    function write_aliased_open_generic_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_aliased_open_generic_written), "Unexpected call to write_aliased_open_generic_");
      expected = obj.write_aliased_open_generic_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_aliased_open_generic_");
      obj.write_aliased_open_generic_written = obj.write_aliased_open_generic_written(2:end);
    end

    function write_aliased_closed_generic_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_aliased_closed_generic_written), "Unexpected call to write_aliased_closed_generic_");
      expected = obj.write_aliased_closed_generic_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_aliased_closed_generic_");
      obj.write_aliased_closed_generic_written = obj.write_aliased_closed_generic_written(2:end);
    end

    function write_aliased_optional_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_aliased_optional_written), "Unexpected call to write_aliased_optional_");
      expected = obj.write_aliased_optional_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_aliased_optional_");
      obj.write_aliased_optional_written = obj.write_aliased_optional_written(2:end);
    end

    function write_aliased_generic_optional_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_aliased_generic_optional_written), "Unexpected call to write_aliased_generic_optional_");
      expected = obj.write_aliased_generic_optional_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_aliased_generic_optional_");
      obj.write_aliased_generic_optional_written = obj.write_aliased_generic_optional_written(2:end);
    end

    function write_aliased_generic_union_2_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_aliased_generic_union_2_written), "Unexpected call to write_aliased_generic_union_2_");
      expected = obj.write_aliased_generic_union_2_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_aliased_generic_union_2_");
      obj.write_aliased_generic_union_2_written = obj.write_aliased_generic_union_2_written(2:end);
    end

    function write_aliased_generic_vector_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_aliased_generic_vector_written), "Unexpected call to write_aliased_generic_vector_");
      expected = obj.write_aliased_generic_vector_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_aliased_generic_vector_");
      obj.write_aliased_generic_vector_written = obj.write_aliased_generic_vector_written(2:end);
    end

    function write_aliased_generic_fixed_vector_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_aliased_generic_fixed_vector_written), "Unexpected call to write_aliased_generic_fixed_vector_");
      expected = obj.write_aliased_generic_fixed_vector_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_aliased_generic_fixed_vector_");
      obj.write_aliased_generic_fixed_vector_written = obj.write_aliased_generic_fixed_vector_written(2:end);
    end

    function write_stream_of_aliased_generic_union_2_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_stream_of_aliased_generic_union_2_written), "Unexpected call to write_stream_of_aliased_generic_union_2_");
      expected = obj.write_stream_of_aliased_generic_union_2_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_stream_of_aliased_generic_union_2_");
      obj.write_stream_of_aliased_generic_union_2_written = obj.write_stream_of_aliased_generic_union_2_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
