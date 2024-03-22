classdef MockFixedArraysWriter < test_model.FixedArraysWriterBase
  properties
    testCase_
    write_ints_written
    write_fixed_simple_record_array_written
    write_fixed_record_with_vlens_array_written
    write_record_with_fixed_arrays_written
    write_named_array_written
  end

  methods
    function obj = MockFixedArraysWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_ints_written = Node.empty();
      obj.write_fixed_simple_record_array_written = Node.empty();
      obj.write_fixed_record_with_vlens_array_written = Node.empty();
      obj.write_record_with_fixed_arrays_written = Node.empty();
      obj.write_named_array_written = Node.empty();
    end

    function expect_write_ints_(obj, value)
      obj.write_ints_written(end+1) = Node(value);
    end

    function expect_write_fixed_simple_record_array_(obj, value)
      obj.write_fixed_simple_record_array_written(end+1) = Node(value);
    end

    function expect_write_fixed_record_with_vlens_array_(obj, value)
      obj.write_fixed_record_with_vlens_array_written(end+1) = Node(value);
    end

    function expect_write_record_with_fixed_arrays_(obj, value)
      obj.write_record_with_fixed_arrays_written(end+1) = Node(value);
    end

    function expect_write_named_array_(obj, value)
      obj.write_named_array_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_ints_written), "Expected call to write_ints_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_fixed_simple_record_array_written), "Expected call to write_fixed_simple_record_array_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_fixed_record_with_vlens_array_written), "Expected call to write_fixed_record_with_vlens_array_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_record_with_fixed_arrays_written), "Expected call to write_record_with_fixed_arrays_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_named_array_written), "Expected call to write_named_array_ was not received");
    end
  end

  methods (Access=protected)
    function write_ints_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_ints_written), "Unexpected call to write_ints_");
      expected = obj.write_ints_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_ints_");
      obj.write_ints_written = obj.write_ints_written(2:end);
    end

    function write_fixed_simple_record_array_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_fixed_simple_record_array_written), "Unexpected call to write_fixed_simple_record_array_");
      expected = obj.write_fixed_simple_record_array_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_fixed_simple_record_array_");
      obj.write_fixed_simple_record_array_written = obj.write_fixed_simple_record_array_written(2:end);
    end

    function write_fixed_record_with_vlens_array_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_fixed_record_with_vlens_array_written), "Unexpected call to write_fixed_record_with_vlens_array_");
      expected = obj.write_fixed_record_with_vlens_array_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_fixed_record_with_vlens_array_");
      obj.write_fixed_record_with_vlens_array_written = obj.write_fixed_record_with_vlens_array_written(2:end);
    end

    function write_record_with_fixed_arrays_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_record_with_fixed_arrays_written), "Unexpected call to write_record_with_fixed_arrays_");
      expected = obj.write_record_with_fixed_arrays_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_record_with_fixed_arrays_");
      obj.write_record_with_fixed_arrays_written = obj.write_record_with_fixed_arrays_written(2:end);
    end

    function write_named_array_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_named_array_written), "Unexpected call to write_named_array_");
      expected = obj.write_named_array_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_named_array_");
      obj.write_named_array_written = obj.write_named_array_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end