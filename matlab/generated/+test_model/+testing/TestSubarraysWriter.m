% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef TestSubarraysWriter < test_model.SubarraysWriterBase
  properties (Access = private)
    writer_
    create_reader_
    mock_writer_
    close_called_
    filename_
    format_
  end

  methods
    function self = TestSubarraysWriter(testCase, format, create_writer, create_reader)
      self.filename_ = tempname();
      self.format_ = format;
      self.writer_ = create_writer(self.filename_);
      self.create_reader_ = create_reader;
      self.mock_writer_ = test_model.testing.MockSubarraysWriter(testCase);
      self.close_called_ = false;
    end

    function delete(self)
      delete(self.filename_);
      if ~self.close_called_
        % ADD_FAILURE() << ...;
        throw(yardl.RuntimeError("Close() must be called on 'TestSubarraysWriter' to verify mocks"));
      end
    end
  end

  methods (Access=protected)
    function write_dynamic_with_fixed_int_subarray_(self, value)
      self.writer_.write_dynamic_with_fixed_int_subarray(value);
      self.mock_writer_.expect_write_dynamic_with_fixed_int_subarray_(value);
    end

    function write_dynamic_with_fixed_float_subarray_(self, value)
      self.writer_.write_dynamic_with_fixed_float_subarray(value);
      self.mock_writer_.expect_write_dynamic_with_fixed_float_subarray_(value);
    end

    function write_known_dim_count_with_fixed_int_subarray_(self, value)
      self.writer_.write_known_dim_count_with_fixed_int_subarray(value);
      self.mock_writer_.expect_write_known_dim_count_with_fixed_int_subarray_(value);
    end

    function write_known_dim_count_with_fixed_float_subarray_(self, value)
      self.writer_.write_known_dim_count_with_fixed_float_subarray(value);
      self.mock_writer_.expect_write_known_dim_count_with_fixed_float_subarray_(value);
    end

    function write_fixed_with_fixed_int_subarray_(self, value)
      self.writer_.write_fixed_with_fixed_int_subarray(value);
      self.mock_writer_.expect_write_fixed_with_fixed_int_subarray_(value);
    end

    function write_fixed_with_fixed_float_subarray_(self, value)
      self.writer_.write_fixed_with_fixed_float_subarray(value);
      self.mock_writer_.expect_write_fixed_with_fixed_float_subarray_(value);
    end

    function write_nested_subarray_(self, value)
      self.writer_.write_nested_subarray(value);
      self.mock_writer_.expect_write_nested_subarray_(value);
    end

    function write_dynamic_with_fixed_vector_subarray_(self, value)
      self.writer_.write_dynamic_with_fixed_vector_subarray(value);
      self.mock_writer_.expect_write_dynamic_with_fixed_vector_subarray_(value);
    end

    function write_generic_subarray_(self, value)
      self.writer_.write_generic_subarray(value);
      self.mock_writer_.expect_write_generic_subarray_(value);
    end

    function close_(self)
      self.close_called_ = true;
      self.writer_.close();
      mock_copy = copy(self.mock_writer_);

      reader = self.create_reader_(self.filename_);
      reader.copy_to(self.mock_writer_);
      reader.close();
      self.mock_writer_.verify();
      self.mock_writer_.close();

      translated = invoke_translator(self.filename_, self.format_, self.format_);
      reader = self.create_reader_(translated);
      reader.copy_to(mock_copy);
      reader.close();
      mock_copy.verify();
      mock_copy.close();
      delete(translated);
    end

    function end_stream_(self)
    end
  end
end
