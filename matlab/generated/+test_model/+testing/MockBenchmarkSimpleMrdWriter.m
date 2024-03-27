% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockBenchmarkSimpleMrdWriter < test_model.BenchmarkSimpleMrdWriterBase
  properties
    testCase_
    write_data_written
  end

  methods
    function obj = MockBenchmarkSimpleMrdWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_data_written = Node.empty();
    end

    function expect_write_data_(obj, value)
      obj.write_data_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_data_written), "Expected call to write_data_ was not received");
    end
  end

  methods (Access=protected)
    function write_data_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_data_written), "Unexpected call to write_data_");
      expected = obj.write_data_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_data_");
      obj.write_data_written = obj.write_data_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
