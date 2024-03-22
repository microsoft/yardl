classdef MockBenchmarkSmallRecordWriter < test_model.BenchmarkSmallRecordWriterBase
  properties
    testCase_
    write_small_record_written
  end

  methods
    function obj = MockBenchmarkSmallRecordWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_small_record_written = Node.empty();
    end

    function expect_write_small_record_(obj, value)
      obj.write_small_record_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_small_record_written), "Expected call to write_small_record_ was not received");
    end
  end

  methods (Access=protected)
    function write_small_record_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_small_record_written), "Unexpected call to write_small_record_");
      expected = obj.write_small_record_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_small_record_");
      obj.write_small_record_written = obj.write_small_record_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end