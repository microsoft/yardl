% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockBenchmarkFloat256x256Writer < matlab.mixin.Copyable & test_model.BenchmarkFloat256x256WriterBase
  properties
    testCase_
    write_float256x256_written
  end

  methods
    function obj = MockBenchmarkFloat256x256Writer(testCase)
      obj.testCase_ = testCase;
      obj.write_float256x256_written = yardl.None;
    end

    function expect_write_float256x256_(obj, value)
      if obj.write_float256x256_written.has_value()
        last_dim = ndims(value);
        obj.write_float256x256_written = yardl.Optional(cat(last_dim, obj.write_float256x256_written.value, value));
      else
        obj.write_float256x256_written = yardl.Optional(value);
      end
    end

    function verify(obj)
      obj.testCase_.verifyEqual(obj.write_float256x256_written, yardl.None, "Expected call to write_float256x256_ was not received");
    end
  end

  methods (Access=protected)
    function write_float256x256_(obj, value)
      obj.testCase_.verifyTrue(obj.write_float256x256_written.has_value(), "Unexpected call to write_float256x256_");
      expected = obj.write_float256x256_written.value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_float256x256_");
      obj.write_float256x256_written = yardl.None;
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
