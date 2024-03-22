classdef TestFixedArraysWriter < test_model.FixedArraysWriterBase
  properties (Access = private)
    writer_
    create_reader_
    mock_writer_
    close_called_
  end

  methods
    function obj = TestFixedArraysWriter(testCase, writer, create_reader)
      obj.writer_ = writer;
      obj.create_reader_ = create_reader;
      obj.mock_writer_ = test_model.testing.MockFixedArraysWriter(testCase);
      obj.close_called_ = false;
    end

    function delete(obj)
      if ~obj.close_called_
        % ADD_FAILURE() << ...;
        throw(yardl.RuntimeError("Close() must be called on 'TestFixedArraysWriter' to verify mocks"));
      end
    end
  end

  methods (Access=protected)
    function write_ints_(obj, value)
      obj.writer_.write_ints(value);
      obj.mock_writer_.expect_write_ints_(value);
    end

    function write_fixed_simple_record_array_(obj, value)
      obj.writer_.write_fixed_simple_record_array(value);
      obj.mock_writer_.expect_write_fixed_simple_record_array_(value);
    end

    function write_fixed_record_with_vlens_array_(obj, value)
      obj.writer_.write_fixed_record_with_vlens_array(value);
      obj.mock_writer_.expect_write_fixed_record_with_vlens_array_(value);
    end

    function write_record_with_fixed_arrays_(obj, value)
      obj.writer_.write_record_with_fixed_arrays(value);
      obj.mock_writer_.expect_write_record_with_fixed_arrays_(value);
    end

    function write_named_array_(obj, value)
      obj.writer_.write_named_array(value);
      obj.mock_writer_.expect_write_named_array_(value);
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
