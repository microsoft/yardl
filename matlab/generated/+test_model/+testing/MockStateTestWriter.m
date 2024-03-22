classdef MockStateTestWriter < test_model.StateTestWriterBase
  properties
    testCase_
    write_an_int_written
    write_a_stream_written
    write_another_int_written
  end

  methods
    function obj = MockStateTestWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_an_int_written = Node.empty();
      obj.write_a_stream_written = Node.empty();
      obj.write_another_int_written = Node.empty();
    end

    function expect_write_an_int_(obj, value)
      obj.write_an_int_written(end+1) = Node(value);
    end

    function expect_write_a_stream_(obj, value)
      obj.write_a_stream_written(end+1) = Node(value);
    end

    function expect_write_another_int_(obj, value)
      obj.write_another_int_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_an_int_written), "Expected call to write_an_int_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_a_stream_written), "Expected call to write_a_stream_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_another_int_written), "Expected call to write_another_int_ was not received");
    end
  end

  methods (Access=protected)
    function write_an_int_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_an_int_written), "Unexpected call to write_an_int_");
      expected = obj.write_an_int_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_an_int_");
      obj.write_an_int_written = obj.write_an_int_written(2:end);
    end

    function write_a_stream_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_a_stream_written), "Unexpected call to write_a_stream_");
      expected = obj.write_a_stream_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_a_stream_");
      obj.write_a_stream_written = obj.write_a_stream_written(2:end);
    end

    function write_another_int_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_another_int_written), "Unexpected call to write_another_int_");
      expected = obj.write_another_int_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_another_int_");
      obj.write_another_int_written = obj.write_another_int_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
