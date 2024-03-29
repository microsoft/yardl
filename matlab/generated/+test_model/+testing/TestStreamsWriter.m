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
    function obj = TestStreamsWriter(testCase, format, create_writer, create_reader)
      obj.filename_ = tempname();
      obj.format_ = format;
      obj.writer_ = create_writer(obj.filename_);
      obj.create_reader_ = create_reader;
      obj.mock_writer_ = test_model.testing.MockStreamsWriter(testCase);
      obj.close_called_ = false;
    end

    function delete(obj)
      % delete(obj.filename_);
      if ~obj.close_called_
        % ADD_FAILURE() << ...;
        throw(yardl.RuntimeError("Close() must be called on 'TestStreamsWriter' to verify mocks"));
      end
    end
  end

  methods (Access=protected)
    function write_int_data_(obj, value)
      obj.writer_.write_int_data(value);
      obj.mock_writer_.expect_write_int_data_(value);
    end

    function write_optional_int_data_(obj, value)
      obj.writer_.write_optional_int_data(value);
      obj.mock_writer_.expect_write_optional_int_data_(value);
    end

    function write_record_with_optional_vector_data_(obj, value)
      obj.writer_.write_record_with_optional_vector_data(value);
      obj.mock_writer_.expect_write_record_with_optional_vector_data_(value);
    end

    function write_fixed_vector_(obj, value)
      obj.writer_.write_fixed_vector(value);
      obj.mock_writer_.expect_write_fixed_vector_(value);
    end

    function close_(obj)
      obj.close_called_ = true;
      obj.writer_.close();
      mock_copy = copy(obj.mock_writer_);

      reader = obj.create_reader_(obj.filename_);
      reader.copy_to(obj.mock_writer_);
      reader.close();
      obj.mock_writer_.verify();
      obj.mock_writer_.close();

      translated = invoke_translator(obj.filename_, obj.format_, obj.format_);
      reader = obj.create_reader_(translated);
      reader.copy_to(mock_copy);
      reader.close();
      mock_copy.verify();
      mock_copy.close();
      % delete(translated);
    end

    function end_stream_(obj)
    end
  end
end
