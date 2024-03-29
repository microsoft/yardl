% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockNDArraysSingleDimensionWriter < matlab.mixin.Copyable & test_model.NDArraysSingleDimensionWriterBase
  properties
    testCase_
    write_ints_written
    write_simple_record_array_written
    write_record_with_vlens_array_written
    write_record_with_nd_arrays_written
  end

  methods
    function obj = MockNDArraysSingleDimensionWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_ints_written = Node.empty();
      obj.write_simple_record_array_written = Node.empty();
      obj.write_record_with_vlens_array_written = Node.empty();
      obj.write_record_with_nd_arrays_written = Node.empty();
    end

    function expect_write_ints_(obj, value)
      if isempty(obj.write_ints_written)
        obj.write_ints_written = Node(value);
      else
        last_dim = ndims(value);
        obj.write_ints_written = Node(cat(last_dim, obj.write_ints_written(1).value, value));
      end
    end

    function expect_write_simple_record_array_(obj, value)
      if isempty(obj.write_simple_record_array_written)
        obj.write_simple_record_array_written = Node(value);
      else
        last_dim = ndims(value);
        obj.write_simple_record_array_written = Node(cat(last_dim, obj.write_simple_record_array_written(1).value, value));
      end
    end

    function expect_write_record_with_vlens_array_(obj, value)
      if isempty(obj.write_record_with_vlens_array_written)
        obj.write_record_with_vlens_array_written = Node(value);
      else
        last_dim = ndims(value);
        obj.write_record_with_vlens_array_written = Node(cat(last_dim, obj.write_record_with_vlens_array_written(1).value, value));
      end
    end

    function expect_write_record_with_nd_arrays_(obj, value)
      if isempty(obj.write_record_with_nd_arrays_written)
        obj.write_record_with_nd_arrays_written = Node(value);
      else
        last_dim = ndims(value);
        obj.write_record_with_nd_arrays_written = Node(cat(last_dim, obj.write_record_with_nd_arrays_written(1).value, value));
      end
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_ints_written), "Expected call to write_ints_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_simple_record_array_written), "Expected call to write_simple_record_array_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_record_with_vlens_array_written), "Expected call to write_record_with_vlens_array_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_record_with_nd_arrays_written), "Expected call to write_record_with_nd_arrays_ was not received");
    end
  end

  methods (Access=protected)
    function write_ints_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_ints_written), "Unexpected call to write_ints_");
      expected = obj.write_ints_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_ints_");
      obj.write_ints_written = Node.empty();
    end

    function write_simple_record_array_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_simple_record_array_written), "Unexpected call to write_simple_record_array_");
      expected = obj.write_simple_record_array_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_simple_record_array_");
      obj.write_simple_record_array_written = Node.empty();
    end

    function write_record_with_vlens_array_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_record_with_vlens_array_written), "Unexpected call to write_record_with_vlens_array_");
      expected = obj.write_record_with_vlens_array_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_record_with_vlens_array_");
      obj.write_record_with_vlens_array_written = Node.empty();
    end

    function write_record_with_nd_arrays_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_record_with_nd_arrays_written), "Unexpected call to write_record_with_nd_arrays_");
      expected = obj.write_record_with_nd_arrays_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_record_with_nd_arrays_");
      obj.write_record_with_nd_arrays_written = Node.empty();
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
