classdef MockMapsWriter < test_model.MapsWriterBase
  properties
    testCase_
    write_string_to_int_written
    write_int_to_string_written
    write_string_to_union_written
    write_aliased_generic_written
  end

  methods
    function obj = MockMapsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_string_to_int_written = Node.empty();
      obj.write_int_to_string_written = Node.empty();
      obj.write_string_to_union_written = Node.empty();
      obj.write_aliased_generic_written = Node.empty();
    end

    function expect_write_string_to_int_(obj, value)
      obj.write_string_to_int_written(end+1) = Node(value);
    end

    function expect_write_int_to_string_(obj, value)
      obj.write_int_to_string_written(end+1) = Node(value);
    end

    function expect_write_string_to_union_(obj, value)
      obj.write_string_to_union_written(end+1) = Node(value);
    end

    function expect_write_aliased_generic_(obj, value)
      obj.write_aliased_generic_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_string_to_int_written), "Expected call to write_string_to_int_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_int_to_string_written), "Expected call to write_int_to_string_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_string_to_union_written), "Expected call to write_string_to_union_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_aliased_generic_written), "Expected call to write_aliased_generic_ was not received");
    end
  end

  methods (Access=protected)
    function write_string_to_int_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_string_to_int_written), "Unexpected call to write_string_to_int_");
      expected = obj.write_string_to_int_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_string_to_int_");
      obj.write_string_to_int_written = obj.write_string_to_int_written(2:end);
    end

    function write_int_to_string_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_int_to_string_written), "Unexpected call to write_int_to_string_");
      expected = obj.write_int_to_string_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_int_to_string_");
      obj.write_int_to_string_written = obj.write_int_to_string_written(2:end);
    end

    function write_string_to_union_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_string_to_union_written), "Unexpected call to write_string_to_union_");
      expected = obj.write_string_to_union_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_string_to_union_");
      obj.write_string_to_union_written = obj.write_string_to_union_written(2:end);
    end

    function write_aliased_generic_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_aliased_generic_written), "Unexpected call to write_aliased_generic_");
      expected = obj.write_aliased_generic_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_aliased_generic_");
      obj.write_aliased_generic_written = obj.write_aliased_generic_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
