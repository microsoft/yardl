% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef TestMapsWriter < test_model.MapsWriterBase
  properties (Access = private)
    writer_
    create_reader_
    mock_writer_
    close_called_
    filename_
    format_
  end

  methods
    function self = TestMapsWriter(testCase, format, create_writer, create_reader)
      self.filename_ = tempname();
      self.format_ = format;
      self.writer_ = create_writer(self.filename_);
      self.create_reader_ = create_reader;
      self.mock_writer_ = test_model.testing.MockMapsWriter(testCase);
      self.close_called_ = false;
    end

    function delete(self)
      delete(self.filename_);
      if ~self.close_called_
        % ADD_FAILURE() << ...;
        throw(yardl.RuntimeError("Close() must be called on 'TestMapsWriter' to verify mocks"));
      end
    end
  end

  methods (Access=protected)
    function write_string_to_int_(self, value)
      self.writer_.write_string_to_int(value);
      self.mock_writer_.expect_write_string_to_int_(value);
    end

    function write_int_to_string_(self, value)
      self.writer_.write_int_to_string(value);
      self.mock_writer_.expect_write_int_to_string_(value);
    end

    function write_string_to_union_(self, value)
      self.writer_.write_string_to_union(value);
      self.mock_writer_.expect_write_string_to_union_(value);
    end

    function write_aliased_generic_(self, value)
      self.writer_.write_aliased_generic(value);
      self.mock_writer_.expect_write_aliased_generic_(value);
    end

    function write_records_(self, value)
      self.writer_.write_records(value);
      self.mock_writer_.expect_write_records_(value);
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
