classdef MockScalarOptionalsWriter < test_model.ScalarOptionalsWriterBase
  properties
    testCase_
    write_optional_int_written
    write_optional_record_written
    write_record_with_optional_fields_written
    write_optional_record_with_optional_fields_written
  end

  methods
    function obj = MockScalarOptionalsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_optional_int_written = Node.empty();
      obj.write_optional_record_written = Node.empty();
      obj.write_record_with_optional_fields_written = Node.empty();
      obj.write_optional_record_with_optional_fields_written = Node.empty();
    end

    function expect_write_optional_int_(obj, value)
      obj.write_optional_int_written(end+1) = Node(value);
    end

    function expect_write_optional_record_(obj, value)
      obj.write_optional_record_written(end+1) = Node(value);
    end

    function expect_write_record_with_optional_fields_(obj, value)
      obj.write_record_with_optional_fields_written(end+1) = Node(value);
    end

    function expect_write_optional_record_with_optional_fields_(obj, value)
      obj.write_optional_record_with_optional_fields_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_optional_int_written), "Expected call to write_optional_int_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_optional_record_written), "Expected call to write_optional_record_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_record_with_optional_fields_written), "Expected call to write_record_with_optional_fields_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_optional_record_with_optional_fields_written), "Expected call to write_optional_record_with_optional_fields_ was not received");
    end
  end

  methods (Access=protected)
    function write_optional_int_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_optional_int_written), "Unexpected call to write_optional_int_");
      expected = obj.write_optional_int_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_optional_int_");
      obj.write_optional_int_written = obj.write_optional_int_written(2:end);
    end

    function write_optional_record_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_optional_record_written), "Unexpected call to write_optional_record_");
      expected = obj.write_optional_record_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_optional_record_");
      obj.write_optional_record_written = obj.write_optional_record_written(2:end);
    end

    function write_record_with_optional_fields_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_record_with_optional_fields_written), "Unexpected call to write_record_with_optional_fields_");
      expected = obj.write_record_with_optional_fields_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_record_with_optional_fields_");
      obj.write_record_with_optional_fields_written = obj.write_record_with_optional_fields_written(2:end);
    end

    function write_optional_record_with_optional_fields_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_optional_record_with_optional_fields_written), "Unexpected call to write_optional_record_with_optional_fields_");
      expected = obj.write_optional_record_with_optional_fields_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_optional_record_with_optional_fields_");
      obj.write_optional_record_with_optional_fields_written = obj.write_optional_record_with_optional_fields_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
