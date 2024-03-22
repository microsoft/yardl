classdef MockProtocolWithKeywordStepsWriter < test_model.ProtocolWithKeywordStepsWriterBase
  properties
    testCase_
    write_int_written
    write_float_written
  end

  methods
    function obj = MockProtocolWithKeywordStepsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_int_written = Node.empty();
      obj.write_float_written = Node.empty();
    end

    function expect_write_int_(obj, value)
      obj.write_int_written(end+1) = Node(value);
    end

    function expect_write_float_(obj, value)
      obj.write_float_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_int_written), "Expected call to write_int_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_float_written), "Expected call to write_float_ was not received");
    end
  end

  methods (Access=protected)
    function write_int_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_int_written), "Unexpected call to write_int_");
      expected = obj.write_int_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_int_");
      obj.write_int_written = obj.write_int_written(2:end);
    end

    function write_float_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_float_written), "Unexpected call to write_float_");
      expected = obj.write_float_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_float_");
      obj.write_float_written = obj.write_float_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end