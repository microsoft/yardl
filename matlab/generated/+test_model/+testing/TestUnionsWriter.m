classdef TestUnionsWriter < test_model.UnionsWriterBase
  properties (Access = private)
    writer_
    create_reader_
    mock_writer_
    close_called_
  end

  methods
    function obj = TestUnionsWriter(testCase, writer, create_reader)
      obj.writer_ = writer;
      obj.create_reader_ = create_reader;
      obj.mock_writer_ = test_model.testing.MockUnionsWriter(testCase);
      obj.close_called_ = false;
    end

    function delete(obj)
      if ~obj.close_called_
        % ADD_FAILURE() << ...;
        throw(yardl.RuntimeError("Close() must be called on 'TestUnionsWriter' to verify mocks"));
      end
    end
  end

  methods (Access=protected)
    function write_int_or_simple_record_(obj, value)
      obj.writer_.write_int_or_simple_record(value);
      obj.mock_writer_.expect_write_int_or_simple_record_(value);
    end

    function write_int_or_record_with_vlens_(obj, value)
      obj.writer_.write_int_or_record_with_vlens(value);
      obj.mock_writer_.expect_write_int_or_record_with_vlens_(value);
    end

    function write_monosotate_or_int_or_simple_record_(obj, value)
      obj.writer_.write_monosotate_or_int_or_simple_record(value);
      obj.mock_writer_.expect_write_monosotate_or_int_or_simple_record_(value);
    end

    function write_record_with_unions_(obj, value)
      obj.writer_.write_record_with_unions(value);
      obj.mock_writer_.expect_write_record_with_unions_(value);
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
