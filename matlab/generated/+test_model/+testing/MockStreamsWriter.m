classdef MockStreamsWriter < test_model.StreamsWriterBase
  properties
    testCase_
    write_int_data_written
    write_optional_int_data_written
    write_record_with_optional_vector_data_written
    write_fixed_vector_written
  end

  methods
    function obj = MockStreamsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_int_data_written = Node.empty();
      obj.write_optional_int_data_written = Node.empty();
      obj.write_record_with_optional_vector_data_written = Node.empty();
      obj.write_fixed_vector_written = Node.empty();
    end

    function expect_write_int_data_(obj, value)
      obj.write_int_data_written(end+1) = Node(value);
    end

    function expect_write_optional_int_data_(obj, value)
      obj.write_optional_int_data_written(end+1) = Node(value);
    end

    function expect_write_record_with_optional_vector_data_(obj, value)
      obj.write_record_with_optional_vector_data_written(end+1) = Node(value);
    end

    function expect_write_fixed_vector_(obj, value)
      obj.write_fixed_vector_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_int_data_written), "Expected call to write_int_data_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_optional_int_data_written), "Expected call to write_optional_int_data_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_record_with_optional_vector_data_written), "Expected call to write_record_with_optional_vector_data_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_fixed_vector_written), "Expected call to write_fixed_vector_ was not received");
    end
  end

  methods (Access=protected)
    function write_int_data_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_int_data_written), "Unexpected call to write_int_data_");
      expected = obj.write_int_data_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_int_data_");
      obj.write_int_data_written = obj.write_int_data_written(2:end);
    end

    function write_optional_int_data_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_optional_int_data_written), "Unexpected call to write_optional_int_data_");
      expected = obj.write_optional_int_data_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_optional_int_data_");
      obj.write_optional_int_data_written = obj.write_optional_int_data_written(2:end);
    end

    function write_record_with_optional_vector_data_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_record_with_optional_vector_data_written), "Unexpected call to write_record_with_optional_vector_data_");
      expected = obj.write_record_with_optional_vector_data_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_record_with_optional_vector_data_");
      obj.write_record_with_optional_vector_data_written = obj.write_record_with_optional_vector_data_written(2:end);
    end

    function write_fixed_vector_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_fixed_vector_written), "Unexpected call to write_fixed_vector_");
      expected = obj.write_fixed_vector_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_fixed_vector_");
      obj.write_fixed_vector_written = obj.write_fixed_vector_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
