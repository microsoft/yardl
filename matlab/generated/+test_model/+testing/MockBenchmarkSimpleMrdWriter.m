% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockBenchmarkSimpleMrdWriter < matlab.mixin.Copyable & test_model.BenchmarkSimpleMrdWriterBase
  properties
    testCase_
    write_data_written
  end

  methods
    function obj = MockBenchmarkSimpleMrdWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_data_written = yardl.None;
    end

    function expect_write_data_(obj, value)
      if obj.write_data_written.has_value()
        last_dim = ndims(value);
        obj.write_data_written = yardl.Optional(cat(last_dim, obj.write_data_written.value, value));
      else
        obj.write_data_written = yardl.Optional(value);
      end
    end

    function verify(obj)
      obj.testCase_.verifyEqual(obj.write_data_written, yardl.None, "Expected call to write_data_ was not received");
    end
  end

  methods (Access=protected)
    function write_data_(obj, value)
      obj.testCase_.verifyTrue(obj.write_data_written.has_value(), "Unexpected call to write_data_");
      expected = obj.write_data_written.value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_data_");
      obj.write_data_written = yardl.None;
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
