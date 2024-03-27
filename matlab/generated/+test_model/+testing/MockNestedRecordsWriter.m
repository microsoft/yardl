% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockNestedRecordsWriter < test_model.NestedRecordsWriterBase
  properties
    testCase_
    write_tuple_with_records_written
  end

  methods
    function obj = MockNestedRecordsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_tuple_with_records_written = Node.empty();
    end

    function expect_write_tuple_with_records_(obj, value)
      obj.write_tuple_with_records_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_tuple_with_records_written), "Expected call to write_tuple_with_records_ was not received");
    end
  end

  methods (Access=protected)
    function write_tuple_with_records_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_tuple_with_records_written), "Unexpected call to write_tuple_with_records_");
      expected = obj.write_tuple_with_records_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_tuple_with_records_");
      obj.write_tuple_with_records_written = obj.write_tuple_with_records_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
