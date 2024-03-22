classdef TestAdvancedGenericsWriter < test_model.AdvancedGenericsWriterBase
  properties (Access = private)
    writer_
    create_reader_
    mock_writer_
    close_called_
  end

  methods
    function obj = TestAdvancedGenericsWriter(testCase, writer, create_reader)
      obj.writer_ = writer;
      obj.create_reader_ = create_reader;
      obj.mock_writer_ = test_model.testing.MockAdvancedGenericsWriter(testCase);
      obj.close_called_ = false;
    end

    function delete(obj)
      if ~obj.close_called_
        % ADD_FAILURE() << ...;
        throw(yardl.RuntimeError("Close() must be called on 'TestAdvancedGenericsWriter' to verify mocks"));
      end
    end
  end

  methods (Access=protected)
    function write_float_image_image_(obj, value)
      obj.writer_.write_float_image_image(value);
      obj.mock_writer_.expect_write_float_image_image_(value);
    end

    function write_generic_record_1_(obj, value)
      obj.writer_.write_generic_record_1(value);
      obj.mock_writer_.expect_write_generic_record_1_(value);
    end

    function write_tuple_of_optionals_(obj, value)
      obj.writer_.write_tuple_of_optionals(value);
      obj.mock_writer_.expect_write_tuple_of_optionals_(value);
    end

    function write_tuple_of_optionals_alternate_syntax_(obj, value)
      obj.writer_.write_tuple_of_optionals_alternate_syntax(value);
      obj.mock_writer_.expect_write_tuple_of_optionals_alternate_syntax_(value);
    end

    function write_tuple_of_vectors_(obj, value)
      obj.writer_.write_tuple_of_vectors(value);
      obj.mock_writer_.expect_write_tuple_of_vectors_(value);
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
