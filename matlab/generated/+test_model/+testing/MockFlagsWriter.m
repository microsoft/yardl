classdef MockFlagsWriter < test_model.FlagsWriterBase
  properties
    testCase_
    write_days_written
    write_formats_written
  end

  methods
    function obj = MockFlagsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_days_written = Node.empty();
      obj.write_formats_written = Node.empty();
    end

    function expect_write_days_(obj, value)
      obj.write_days_written(end+1) = Node(value);
    end

    function expect_write_formats_(obj, value)
      obj.write_formats_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_days_written), "Expected call to write_days_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_formats_written), "Expected call to write_formats_ was not received");
    end
  end

  methods (Access=protected)
    function write_days_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_days_written), "Unexpected call to write_days_");
      expected = obj.write_days_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_days_");
      obj.write_days_written = obj.write_days_written(2:end);
    end

    function write_formats_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_formats_written), "Unexpected call to write_formats_");
      expected = obj.write_formats_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_formats_");
      obj.write_formats_written = obj.write_formats_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
