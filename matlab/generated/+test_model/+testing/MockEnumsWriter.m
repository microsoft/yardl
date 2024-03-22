classdef MockEnumsWriter < test_model.EnumsWriterBase
  properties
    testCase_
    write_single_written
    write_vec_written
    write_size_written
  end

  methods
    function obj = MockEnumsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_single_written = Node.empty();
      obj.write_vec_written = Node.empty();
      obj.write_size_written = Node.empty();
    end

    function expect_write_single_(obj, value)
      obj.write_single_written(end+1) = Node(value);
    end

    function expect_write_vec_(obj, value)
      obj.write_vec_written(end+1) = Node(value);
    end

    function expect_write_size_(obj, value)
      obj.write_size_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_single_written), "Expected call to write_single_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_vec_written), "Expected call to write_vec_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_size_written), "Expected call to write_size_ was not received");
    end
  end

  methods (Access=protected)
    function write_single_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_single_written), "Unexpected call to write_single_");
      expected = obj.write_single_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_single_");
      obj.write_single_written = obj.write_single_written(2:end);
    end

    function write_vec_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_vec_written), "Unexpected call to write_vec_");
      expected = obj.write_vec_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_vec_");
      obj.write_vec_written = obj.write_vec_written(2:end);
    end

    function write_size_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_size_written), "Unexpected call to write_size_");
      expected = obj.write_size_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_size_");
      obj.write_size_written = obj.write_size_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
