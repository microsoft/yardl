% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef TestFlagsWriter < test_model.FlagsWriterBase
  properties (Access = private)
    writer_
    create_reader_
    mock_writer_
    close_called_
    filename_
    format_
  end

  methods
    function obj = TestFlagsWriter(testCase, format, create_writer, create_reader)
      obj.filename_ = tempname();
      obj.format_ = format;
      obj.writer_ = create_writer(obj.filename_);
      obj.create_reader_ = create_reader;
      obj.mock_writer_ = test_model.testing.MockFlagsWriter(testCase);
      obj.close_called_ = false;
    end

    function delete(obj)
      delete(obj.filename_);
      if ~obj.close_called_
        % ADD_FAILURE() << ...;
        throw(yardl.RuntimeError("Close() must be called on 'TestFlagsWriter' to verify mocks"));
      end
    end
    function end_days(obj)
      end_days@test_model.FlagsWriterBase(obj);
      obj.writer_.end_days();
    end

    function end_formats(obj)
      end_formats@test_model.FlagsWriterBase(obj);
      obj.writer_.end_formats();
    end

  end

  methods (Access=protected)
    function write_days_(obj, value)
      obj.writer_.write_days(value);
      obj.mock_writer_.expect_write_days_(value);
    end

    function write_formats_(obj, value)
      obj.writer_.write_formats(value);
      obj.mock_writer_.expect_write_formats_(value);
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
      delete(translated);
    end

    function end_stream_(obj)
    end
  end
end
