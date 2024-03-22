classdef MockOptionalVectorsWriter < test_model.OptionalVectorsWriterBase
  properties
    testCase_
    write_record_with_optional_vector_written
  end

  methods
    function obj = MockOptionalVectorsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_record_with_optional_vector_written = Node.empty();
    end

    function expect_write_record_with_optional_vector_(obj, value)
      obj.write_record_with_optional_vector_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_record_with_optional_vector_written), "Expected call to write_record_with_optional_vector_ was not received");
    end
  end

  methods (Access=protected)
    function write_record_with_optional_vector_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_record_with_optional_vector_written), "Unexpected call to write_record_with_optional_vector_");
      expected = obj.write_record_with_optional_vector_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_record_with_optional_vector_");
      obj.write_record_with_optional_vector_written = obj.write_record_with_optional_vector_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
