classdef TestMapsWriter < test_model.MapsWriterBase
  properties (Access = private)
    writer_
    create_reader_
    mock_writer_
    close_called_
  end

  methods
    function obj = TestMapsWriter(testCase, writer, create_reader)
      obj.writer_ = writer;
      obj.create_reader_ = create_reader;
      obj.mock_writer_ = test_model.testing.MockMapsWriter(testCase);
      obj.close_called_ = false;
    end

    function delete(obj)
      if ~obj.close_called_
        % ADD_FAILURE() << ...;
        throw(yardl.RuntimeError("Close() must be called on 'TestMapsWriter' to verify mocks"));
      end
    end
  end

  methods (Access=protected)
    function write_string_to_int_(obj, value)
      obj.writer_.write_string_to_int(value);
      obj.mock_writer_.expect_write_string_to_int_(value);
    end

    function write_int_to_string_(obj, value)
      obj.writer_.write_int_to_string(value);
      obj.mock_writer_.expect_write_int_to_string_(value);
    end

    function write_string_to_union_(obj, value)
      obj.writer_.write_string_to_union(value);
      obj.mock_writer_.expect_write_string_to_union_(value);
    end

    function write_aliased_generic_(obj, value)
      obj.writer_.write_aliased_generic(value);
      obj.mock_writer_.expect_write_aliased_generic_(value);
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
