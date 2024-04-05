% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockNestedRecordsWriter < matlab.mixin.Copyable & test_model.NestedRecordsWriterBase
  properties
    testCase_
    write_tuple_with_records_written
  end

  methods
    function obj = MockNestedRecordsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_tuple_with_records_written = yardl.None;
    end

    function expect_write_tuple_with_records_(obj, value)
      if obj.write_tuple_with_records_written.has_value()
        last_dim = ndims(value);
        obj.write_tuple_with_records_written = yardl.Optional(cat(last_dim, obj.write_tuple_with_records_written.value, value));
      else
        obj.write_tuple_with_records_written = yardl.Optional(value);
      end
    end

    function verify(obj)
      obj.testCase_.verifyEqual(obj.write_tuple_with_records_written, yardl.None, "Expected call to write_tuple_with_records_ was not received");
    end
  end

  methods (Access=protected)
    function write_tuple_with_records_(obj, value)
      obj.testCase_.verifyTrue(obj.write_tuple_with_records_written.has_value(), "Unexpected call to write_tuple_with_records_");
      expected = obj.write_tuple_with_records_written.value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_tuple_with_records_");
      obj.write_tuple_with_records_written = yardl.None;
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
