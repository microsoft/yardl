classdef MockProtocolWithComputedFieldsWriter < test_model.ProtocolWithComputedFieldsWriterBase
  properties
    testCase_
    write_record_with_computed_fields_written
  end

  methods
    function obj = MockProtocolWithComputedFieldsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_record_with_computed_fields_written = Node.empty();
    end

    function expect_write_record_with_computed_fields_(obj, value)
      obj.write_record_with_computed_fields_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_record_with_computed_fields_written), "Expected call to write_record_with_computed_fields_ was not received");
    end
  end

  methods (Access=protected)
    function write_record_with_computed_fields_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_record_with_computed_fields_written), "Unexpected call to write_record_with_computed_fields_");
      expected = obj.write_record_with_computed_fields_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_record_with_computed_fields_");
      obj.write_record_with_computed_fields_written = obj.write_record_with_computed_fields_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end