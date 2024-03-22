classdef TestSubarraysWriter < test_model.SubarraysWriterBase
  properties (Access = private)
    writer_
    create_reader_
    mock_writer_
    close_called_
  end

  methods
    function obj = TestSubarraysWriter(testCase, writer, create_reader)
      obj.writer_ = writer;
      obj.create_reader_ = create_reader;
      obj.mock_writer_ = test_model.testing.MockSubarraysWriter(testCase);
      obj.close_called_ = false;
    end

    function delete(obj)
      if ~obj.close_called_
        % ADD_FAILURE() << ...;
        throw(yardl.RuntimeError("Close() must be called on 'TestSubarraysWriter' to verify mocks"));
      end
    end
  end

  methods (Access=protected)
    function write_dynamic_with_fixed_int_subarray_(obj, value)
      obj.writer_.write_dynamic_with_fixed_int_subarray(value);
      obj.mock_writer_.expect_write_dynamic_with_fixed_int_subarray_(value);
    end

    function write_dynamic_with_fixed_float_subarray_(obj, value)
      obj.writer_.write_dynamic_with_fixed_float_subarray(value);
      obj.mock_writer_.expect_write_dynamic_with_fixed_float_subarray_(value);
    end

    function write_known_dim_count_with_fixed_int_subarray_(obj, value)
      obj.writer_.write_known_dim_count_with_fixed_int_subarray(value);
      obj.mock_writer_.expect_write_known_dim_count_with_fixed_int_subarray_(value);
    end

    function write_known_dim_count_with_fixed_float_subarray_(obj, value)
      obj.writer_.write_known_dim_count_with_fixed_float_subarray(value);
      obj.mock_writer_.expect_write_known_dim_count_with_fixed_float_subarray_(value);
    end

    function write_fixed_with_fixed_int_subarray_(obj, value)
      obj.writer_.write_fixed_with_fixed_int_subarray(value);
      obj.mock_writer_.expect_write_fixed_with_fixed_int_subarray_(value);
    end

    function write_fixed_with_fixed_float_subarray_(obj, value)
      obj.writer_.write_fixed_with_fixed_float_subarray(value);
      obj.mock_writer_.expect_write_fixed_with_fixed_float_subarray_(value);
    end

    function write_nested_subarray_(obj, value)
      obj.writer_.write_nested_subarray(value);
      obj.mock_writer_.expect_write_nested_subarray_(value);
    end

    function write_dynamic_with_fixed_vector_subarray_(obj, value)
      obj.writer_.write_dynamic_with_fixed_vector_subarray(value);
      obj.mock_writer_.expect_write_dynamic_with_fixed_vector_subarray_(value);
    end

    function write_generic_subarray_(obj, value)
      obj.writer_.write_generic_subarray(value);
      obj.mock_writer_.expect_write_generic_subarray_(value);
    end

    function close_(obj)
      obj.close_called_ = true;
      obj.writer_.close();
      reader = obj.create_reader_();
      reader.copy_to(obj.mock_writer_);
      reader.close();
      obj.mock_writer_.verify();
    end

    function end_stream_(obj)
    end
  end
end
