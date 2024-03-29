% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockBenchmarkInt256x256Writer < matlab.mixin.Copyable & test_model.BenchmarkInt256x256WriterBase
  properties
    testCase_
    write_int256x256_written
  end

  methods
    function obj = MockBenchmarkInt256x256Writer(testCase)
      obj.testCase_ = testCase;
      obj.write_int256x256_written = Node.empty();
    end

    function expect_write_int256x256_(obj, value)
      if isempty(obj.write_int256x256_written)
        obj.write_int256x256_written = Node(value);
      else
        last_dim = ndims(value);
        obj.write_int256x256_written = Node(cat(last_dim, obj.write_int256x256_written(1).value, value));
      end
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_int256x256_written), "Expected call to write_int256x256_ was not received");
    end
  end

  methods (Access=protected)
    function write_int256x256_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_int256x256_written), "Unexpected call to write_int256x256_");
      expected = obj.write_int256x256_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_int256x256_");
      obj.write_int256x256_written = Node.empty();
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
