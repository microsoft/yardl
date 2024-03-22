classdef MockStringsWriter < test_model.StringsWriterBase
  properties
    testCase_
    write_single_string_written
    write_rec_with_string_written
  end

  methods
    function obj = MockStringsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_single_string_written = Node.empty();
      obj.write_rec_with_string_written = Node.empty();
    end

    function expect_write_single_string_(obj, value)
      obj.write_single_string_written(end+1) = Node(value);
    end

    function expect_write_rec_with_string_(obj, value)
      obj.write_rec_with_string_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_single_string_written), "Expected call to write_single_string_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_rec_with_string_written), "Expected call to write_rec_with_string_ was not received");
    end
  end

  methods (Access=protected)
    function write_single_string_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_single_string_written), "Unexpected call to write_single_string_");
      expected = obj.write_single_string_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_single_string_");
      obj.write_single_string_written = obj.write_single_string_written(2:end);
    end

    function write_rec_with_string_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_rec_with_string_written), "Unexpected call to write_rec_with_string_");
      expected = obj.write_rec_with_string_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_rec_with_string_");
      obj.write_rec_with_string_written = obj.write_rec_with_string_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
