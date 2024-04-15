% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef TestStreamsWriter < test_model.StreamsWriterBase
  properties (Access = private)
    writer_
    create_reader_
    mock_writer_
    close_called_
    filename_
    format_
  end

  methods
    function self = TestStreamsWriter(testCase, format, create_writer, create_reader)
      self.filename_ = tempname();
      self.format_ = format;
      self.writer_ = create_writer(self.filename_);
      self.create_reader_ = create_reader;
      self.mock_writer_ = test_model.testing.MockStreamsWriter(testCase);
      self.close_called_ = false;
    end

    function delete(self)
      delete(self.filename_);
      if ~self.close_called_
        % ADD_FAILURE() << ...;
        throw(yardl.RuntimeError("Close() must be called on 'TestStreamsWriter' to verify mocks"));
      end
    end
    function end_int_data(self)
      end_int_data@test_model.StreamsWriterBase(self);
      self.writer_.end_int_data();
    end

    function end_optional_int_data(self)
      end_optional_int_data@test_model.StreamsWriterBase(self);
      self.writer_.end_optional_int_data();
    end

    function end_record_with_optional_vector_data(self)
      end_record_with_optional_vector_data@test_model.StreamsWriterBase(self);
      self.writer_.end_record_with_optional_vector_data();
    end

    function end_fixed_vector(self)
      end_fixed_vector@test_model.StreamsWriterBase(self);
      self.writer_.end_fixed_vector();
    end

  end

  methods (Access=protected)
    function write_int_data_(self, value)
      self.writer_.write_int_data(value);
      self.mock_writer_.expect_write_int_data_(value);
    end

    function write_optional_int_data_(self, value)
      self.writer_.write_optional_int_data(value);
      self.mock_writer_.expect_write_optional_int_data_(value);
    end

    function write_record_with_optional_vector_data_(self, value)
      self.writer_.write_record_with_optional_vector_data(value);
      self.mock_writer_.expect_write_record_with_optional_vector_data_(value);
    end

    function write_fixed_vector_(self, value)
      self.writer_.write_fixed_vector(value);
      self.mock_writer_.expect_write_fixed_vector_(value);
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
