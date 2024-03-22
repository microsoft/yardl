classdef TestSimpleGenericsWriter < test_model.SimpleGenericsWriterBase
  properties (Access = private)
    writer_
    create_reader_
    mock_writer_
    close_called_
  end

  methods
    function obj = TestSimpleGenericsWriter(testCase, writer, create_reader)
      obj.writer_ = writer;
      obj.create_reader_ = create_reader;
      obj.mock_writer_ = test_model.testing.MockSimpleGenericsWriter(testCase);
      obj.close_called_ = false;
    end

    function delete(obj)
      if ~obj.close_called_
        % ADD_FAILURE() << ...;
        throw(yardl.RuntimeError("Close() must be called on 'TestSimpleGenericsWriter' to verify mocks"));
      end
    end
  end

  methods (Access=protected)
    function write_float_image_(obj, value)
      obj.writer_.write_float_image(value);
      obj.mock_writer_.expect_write_float_image_(value);
    end

    function write_int_image_(obj, value)
      obj.writer_.write_int_image(value);
      obj.mock_writer_.expect_write_int_image_(value);
    end

    function write_int_image_alternate_syntax_(obj, value)
      obj.writer_.write_int_image_alternate_syntax(value);
      obj.mock_writer_.expect_write_int_image_alternate_syntax_(value);
    end

    function write_string_image_(obj, value)
      obj.writer_.write_string_image(value);
      obj.mock_writer_.expect_write_string_image_(value);
    end

    function write_int_float_tuple_(obj, value)
      obj.writer_.write_int_float_tuple(value);
      obj.mock_writer_.expect_write_int_float_tuple_(value);
    end

    function write_float_float_tuple_(obj, value)
      obj.writer_.write_float_float_tuple(value);
      obj.mock_writer_.expect_write_float_float_tuple_(value);
    end

    function write_int_float_tuple_alternate_syntax_(obj, value)
      obj.writer_.write_int_float_tuple_alternate_syntax(value);
      obj.mock_writer_.expect_write_int_float_tuple_alternate_syntax_(value);
    end

    function write_int_string_tuple_(obj, value)
      obj.writer_.write_int_string_tuple(value);
      obj.mock_writer_.expect_write_int_string_tuple_(value);
    end

    function write_stream_of_type_variants_(obj, value)
      obj.writer_.write_stream_of_type_variants(value);
      obj.mock_writer_.expect_write_stream_of_type_variants_(value);
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
