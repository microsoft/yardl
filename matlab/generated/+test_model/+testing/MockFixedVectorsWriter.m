% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockFixedVectorsWriter < matlab.mixin.Copyable & test_model.FixedVectorsWriterBase
  properties
    testCase_
    write_fixed_int_vector_written
    write_fixed_simple_record_vector_written
    write_fixed_record_with_vlens_vector_written
    write_record_with_fixed_vectors_written
  end

  methods
    function obj = MockFixedVectorsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_fixed_int_vector_written = Node.empty();
      obj.write_fixed_simple_record_vector_written = Node.empty();
      obj.write_fixed_record_with_vlens_vector_written = Node.empty();
      obj.write_record_with_fixed_vectors_written = Node.empty();
    end

    function expect_write_fixed_int_vector_(obj, value)
      if isempty(obj.write_fixed_int_vector_written)
        obj.write_fixed_int_vector_written = Node(value);
      else
        last_dim = ndims(value);
        obj.write_fixed_int_vector_written = Node(cat(last_dim, obj.write_fixed_int_vector_written(1).value, value));
      end
    end

    function expect_write_fixed_simple_record_vector_(obj, value)
      if isempty(obj.write_fixed_simple_record_vector_written)
        obj.write_fixed_simple_record_vector_written = Node(value);
      else
        last_dim = ndims(value);
        obj.write_fixed_simple_record_vector_written = Node(cat(last_dim, obj.write_fixed_simple_record_vector_written(1).value, value));
      end
    end

    function expect_write_fixed_record_with_vlens_vector_(obj, value)
      if isempty(obj.write_fixed_record_with_vlens_vector_written)
        obj.write_fixed_record_with_vlens_vector_written = Node(value);
      else
        last_dim = ndims(value);
        obj.write_fixed_record_with_vlens_vector_written = Node(cat(last_dim, obj.write_fixed_record_with_vlens_vector_written(1).value, value));
      end
    end

    function expect_write_record_with_fixed_vectors_(obj, value)
      if isempty(obj.write_record_with_fixed_vectors_written)
        obj.write_record_with_fixed_vectors_written = Node(value);
      else
        last_dim = ndims(value);
        obj.write_record_with_fixed_vectors_written = Node(cat(last_dim, obj.write_record_with_fixed_vectors_written(1).value, value));
      end
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_fixed_int_vector_written), "Expected call to write_fixed_int_vector_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_fixed_simple_record_vector_written), "Expected call to write_fixed_simple_record_vector_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_fixed_record_with_vlens_vector_written), "Expected call to write_fixed_record_with_vlens_vector_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_record_with_fixed_vectors_written), "Expected call to write_record_with_fixed_vectors_ was not received");
    end
  end

  methods (Access=protected)
    function write_fixed_int_vector_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_fixed_int_vector_written), "Unexpected call to write_fixed_int_vector_");
      expected = obj.write_fixed_int_vector_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_fixed_int_vector_");
      obj.write_fixed_int_vector_written = Node.empty();
    end

    function write_fixed_simple_record_vector_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_fixed_simple_record_vector_written), "Unexpected call to write_fixed_simple_record_vector_");
      expected = obj.write_fixed_simple_record_vector_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_fixed_simple_record_vector_");
      obj.write_fixed_simple_record_vector_written = Node.empty();
    end

    function write_fixed_record_with_vlens_vector_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_fixed_record_with_vlens_vector_written), "Unexpected call to write_fixed_record_with_vlens_vector_");
      expected = obj.write_fixed_record_with_vlens_vector_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_fixed_record_with_vlens_vector_");
      obj.write_fixed_record_with_vlens_vector_written = Node.empty();
    end

    function write_record_with_fixed_vectors_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_record_with_fixed_vectors_written), "Unexpected call to write_record_with_fixed_vectors_");
      expected = obj.write_record_with_fixed_vectors_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_record_with_fixed_vectors_");
      obj.write_record_with_fixed_vectors_written = Node.empty();
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
