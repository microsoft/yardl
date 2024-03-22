classdef MockVlensWriter < test_model.VlensWriterBase
  properties
    testCase_
    write_int_vector_written
    write_complex_vector_written
    write_record_with_vlens_written
    write_vlen_of_record_with_vlens_written
  end

  methods
    function obj = MockVlensWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_int_vector_written = Node.empty();
      obj.write_complex_vector_written = Node.empty();
      obj.write_record_with_vlens_written = Node.empty();
      obj.write_vlen_of_record_with_vlens_written = Node.empty();
    end

    function expect_write_int_vector_(obj, value)
      obj.write_int_vector_written(end+1) = Node(value);
    end

    function expect_write_complex_vector_(obj, value)
      obj.write_complex_vector_written(end+1) = Node(value);
    end

    function expect_write_record_with_vlens_(obj, value)
      obj.write_record_with_vlens_written(end+1) = Node(value);
    end

    function expect_write_vlen_of_record_with_vlens_(obj, value)
      obj.write_vlen_of_record_with_vlens_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_int_vector_written), "Expected call to write_int_vector_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_complex_vector_written), "Expected call to write_complex_vector_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_record_with_vlens_written), "Expected call to write_record_with_vlens_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_vlen_of_record_with_vlens_written), "Expected call to write_vlen_of_record_with_vlens_ was not received");
    end
  end

  methods (Access=protected)
    function write_int_vector_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_int_vector_written), "Unexpected call to write_int_vector_");
      expected = obj.write_int_vector_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_int_vector_");
      obj.write_int_vector_written = obj.write_int_vector_written(2:end);
    end

    function write_complex_vector_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_complex_vector_written), "Unexpected call to write_complex_vector_");
      expected = obj.write_complex_vector_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_complex_vector_");
      obj.write_complex_vector_written = obj.write_complex_vector_written(2:end);
    end

    function write_record_with_vlens_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_record_with_vlens_written), "Unexpected call to write_record_with_vlens_");
      expected = obj.write_record_with_vlens_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_record_with_vlens_");
      obj.write_record_with_vlens_written = obj.write_record_with_vlens_written(2:end);
    end

    function write_vlen_of_record_with_vlens_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_vlen_of_record_with_vlens_written), "Unexpected call to write_vlen_of_record_with_vlens_");
      expected = obj.write_vlen_of_record_with_vlens_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_vlen_of_record_with_vlens_");
      obj.write_vlen_of_record_with_vlens_written = obj.write_vlen_of_record_with_vlens_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
