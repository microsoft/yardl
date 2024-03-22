classdef MockSubarraysWriter < test_model.SubarraysWriterBase
  properties
    testCase_
    write_dynamic_with_fixed_int_subarray_written
    write_dynamic_with_fixed_float_subarray_written
    write_known_dim_count_with_fixed_int_subarray_written
    write_known_dim_count_with_fixed_float_subarray_written
    write_fixed_with_fixed_int_subarray_written
    write_fixed_with_fixed_float_subarray_written
    write_nested_subarray_written
    write_dynamic_with_fixed_vector_subarray_written
    write_generic_subarray_written
  end

  methods
    function obj = MockSubarraysWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_dynamic_with_fixed_int_subarray_written = Node.empty();
      obj.write_dynamic_with_fixed_float_subarray_written = Node.empty();
      obj.write_known_dim_count_with_fixed_int_subarray_written = Node.empty();
      obj.write_known_dim_count_with_fixed_float_subarray_written = Node.empty();
      obj.write_fixed_with_fixed_int_subarray_written = Node.empty();
      obj.write_fixed_with_fixed_float_subarray_written = Node.empty();
      obj.write_nested_subarray_written = Node.empty();
      obj.write_dynamic_with_fixed_vector_subarray_written = Node.empty();
      obj.write_generic_subarray_written = Node.empty();
    end

    function expect_write_dynamic_with_fixed_int_subarray_(obj, value)
      obj.write_dynamic_with_fixed_int_subarray_written(end+1) = Node(value);
    end

    function expect_write_dynamic_with_fixed_float_subarray_(obj, value)
      obj.write_dynamic_with_fixed_float_subarray_written(end+1) = Node(value);
    end

    function expect_write_known_dim_count_with_fixed_int_subarray_(obj, value)
      obj.write_known_dim_count_with_fixed_int_subarray_written(end+1) = Node(value);
    end

    function expect_write_known_dim_count_with_fixed_float_subarray_(obj, value)
      obj.write_known_dim_count_with_fixed_float_subarray_written(end+1) = Node(value);
    end

    function expect_write_fixed_with_fixed_int_subarray_(obj, value)
      obj.write_fixed_with_fixed_int_subarray_written(end+1) = Node(value);
    end

    function expect_write_fixed_with_fixed_float_subarray_(obj, value)
      obj.write_fixed_with_fixed_float_subarray_written(end+1) = Node(value);
    end

    function expect_write_nested_subarray_(obj, value)
      obj.write_nested_subarray_written(end+1) = Node(value);
    end

    function expect_write_dynamic_with_fixed_vector_subarray_(obj, value)
      obj.write_dynamic_with_fixed_vector_subarray_written(end+1) = Node(value);
    end

    function expect_write_generic_subarray_(obj, value)
      obj.write_generic_subarray_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_dynamic_with_fixed_int_subarray_written), "Expected call to write_dynamic_with_fixed_int_subarray_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_dynamic_with_fixed_float_subarray_written), "Expected call to write_dynamic_with_fixed_float_subarray_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_known_dim_count_with_fixed_int_subarray_written), "Expected call to write_known_dim_count_with_fixed_int_subarray_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_known_dim_count_with_fixed_float_subarray_written), "Expected call to write_known_dim_count_with_fixed_float_subarray_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_fixed_with_fixed_int_subarray_written), "Expected call to write_fixed_with_fixed_int_subarray_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_fixed_with_fixed_float_subarray_written), "Expected call to write_fixed_with_fixed_float_subarray_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_nested_subarray_written), "Expected call to write_nested_subarray_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_dynamic_with_fixed_vector_subarray_written), "Expected call to write_dynamic_with_fixed_vector_subarray_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_generic_subarray_written), "Expected call to write_generic_subarray_ was not received");
    end
  end

  methods (Access=protected)
    function write_dynamic_with_fixed_int_subarray_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_dynamic_with_fixed_int_subarray_written), "Unexpected call to write_dynamic_with_fixed_int_subarray_");
      expected = obj.write_dynamic_with_fixed_int_subarray_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_dynamic_with_fixed_int_subarray_");
      obj.write_dynamic_with_fixed_int_subarray_written = obj.write_dynamic_with_fixed_int_subarray_written(2:end);
    end

    function write_dynamic_with_fixed_float_subarray_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_dynamic_with_fixed_float_subarray_written), "Unexpected call to write_dynamic_with_fixed_float_subarray_");
      expected = obj.write_dynamic_with_fixed_float_subarray_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_dynamic_with_fixed_float_subarray_");
      obj.write_dynamic_with_fixed_float_subarray_written = obj.write_dynamic_with_fixed_float_subarray_written(2:end);
    end

    function write_known_dim_count_with_fixed_int_subarray_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_known_dim_count_with_fixed_int_subarray_written), "Unexpected call to write_known_dim_count_with_fixed_int_subarray_");
      expected = obj.write_known_dim_count_with_fixed_int_subarray_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_known_dim_count_with_fixed_int_subarray_");
      obj.write_known_dim_count_with_fixed_int_subarray_written = obj.write_known_dim_count_with_fixed_int_subarray_written(2:end);
    end

    function write_known_dim_count_with_fixed_float_subarray_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_known_dim_count_with_fixed_float_subarray_written), "Unexpected call to write_known_dim_count_with_fixed_float_subarray_");
      expected = obj.write_known_dim_count_with_fixed_float_subarray_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_known_dim_count_with_fixed_float_subarray_");
      obj.write_known_dim_count_with_fixed_float_subarray_written = obj.write_known_dim_count_with_fixed_float_subarray_written(2:end);
    end

    function write_fixed_with_fixed_int_subarray_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_fixed_with_fixed_int_subarray_written), "Unexpected call to write_fixed_with_fixed_int_subarray_");
      expected = obj.write_fixed_with_fixed_int_subarray_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_fixed_with_fixed_int_subarray_");
      obj.write_fixed_with_fixed_int_subarray_written = obj.write_fixed_with_fixed_int_subarray_written(2:end);
    end

    function write_fixed_with_fixed_float_subarray_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_fixed_with_fixed_float_subarray_written), "Unexpected call to write_fixed_with_fixed_float_subarray_");
      expected = obj.write_fixed_with_fixed_float_subarray_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_fixed_with_fixed_float_subarray_");
      obj.write_fixed_with_fixed_float_subarray_written = obj.write_fixed_with_fixed_float_subarray_written(2:end);
    end

    function write_nested_subarray_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_nested_subarray_written), "Unexpected call to write_nested_subarray_");
      expected = obj.write_nested_subarray_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_nested_subarray_");
      obj.write_nested_subarray_written = obj.write_nested_subarray_written(2:end);
    end

    function write_dynamic_with_fixed_vector_subarray_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_dynamic_with_fixed_vector_subarray_written), "Unexpected call to write_dynamic_with_fixed_vector_subarray_");
      expected = obj.write_dynamic_with_fixed_vector_subarray_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_dynamic_with_fixed_vector_subarray_");
      obj.write_dynamic_with_fixed_vector_subarray_written = obj.write_dynamic_with_fixed_vector_subarray_written(2:end);
    end

    function write_generic_subarray_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_generic_subarray_written), "Unexpected call to write_generic_subarray_");
      expected = obj.write_generic_subarray_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_generic_subarray_");
      obj.write_generic_subarray_written = obj.write_generic_subarray_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
