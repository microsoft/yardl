classdef MockBenchmarkFloatVlenWriter < test_model.BenchmarkFloatVlenWriterBase
  properties
    testCase_
    write_float_array_written
  end

  methods
    function obj = MockBenchmarkFloatVlenWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_float_array_written = Node.empty();
    end

    function expect_write_float_array_(obj, value)
      obj.write_float_array_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_float_array_written), "Expected call to write_float_array_ was not received");
    end
  end

  methods (Access=protected)
    function write_float_array_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_float_array_written), "Unexpected call to write_float_array_");
      expected = obj.write_float_array_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_float_array_");
      obj.write_float_array_written = obj.write_float_array_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end