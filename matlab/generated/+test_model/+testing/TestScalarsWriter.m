classdef TestScalarsWriter < test_model.ScalarsWriterBase
  properties (Access = private)
    writer_
    create_reader_
    mock_writer_
    close_called_
  end

  methods
    function obj = TestScalarsWriter(testCase, writer, create_reader)
      obj.writer_ = writer;
      obj.create_reader_ = create_reader;
      obj.mock_writer_ = test_model.testing.MockScalarsWriter(testCase);
      obj.close_called_ = false;
    end

    function delete(obj)
      if ~obj.close_called_
        % ADD_FAILURE() << ...;
        throw(yardl.RuntimeError("Close() must be called on 'TestScalarsWriter' to verify mocks"));
      end
    end
  end

  methods (Access=protected)
    function write_int32_(obj, value)
      obj.writer_.write_int32(value);
      obj.mock_writer_.expect_write_int32_(value);
    end

    function write_record_(obj, value)
      obj.writer_.write_record(value);
      obj.mock_writer_.expect_write_record_(value);
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
