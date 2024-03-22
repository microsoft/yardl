classdef MockBenchmarkFloat256x256Writer < test_model.BenchmarkFloat256x256WriterBase
  properties
    testCase_
    write_float256x256_written
  end

  methods
    function obj = MockBenchmarkFloat256x256Writer(testCase)
      obj.testCase_ = testCase;
      obj.write_float256x256_written = Node.empty();
    end

    function expect_write_float256x256_(obj, value)
      obj.write_float256x256_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_float256x256_written), "Expected call to write_float256x256_ was not received");
    end
  end

  methods (Access=protected)
    function write_float256x256_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_float256x256_written), "Unexpected call to write_float256x256_");
      expected = obj.write_float256x256_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_float256x256_");
      obj.write_float256x256_written = obj.write_float256x256_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
