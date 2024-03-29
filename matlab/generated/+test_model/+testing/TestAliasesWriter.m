% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef TestAliasesWriter < test_model.AliasesWriterBase
  properties (Access = private)
    writer_
    create_reader_
    mock_writer_
    close_called_
    filename_
    format_
  end

  methods
    function obj = TestAliasesWriter(testCase, format, create_writer, create_reader)
      obj.filename_ = tempname();
      obj.format_ = format;
      obj.writer_ = create_writer(obj.filename_);
      obj.create_reader_ = create_reader;
      obj.mock_writer_ = test_model.testing.MockAliasesWriter(testCase);
      obj.close_called_ = false;
    end

    function delete(obj)
      % delete(obj.filename_);
      if ~obj.close_called_
        % ADD_FAILURE() << ...;
        throw(yardl.RuntimeError("Close() must be called on 'TestAliasesWriter' to verify mocks"));
      end
    end
  end

  methods (Access=protected)
    function write_aliased_string_(obj, value)
      obj.writer_.write_aliased_string(value);
      obj.mock_writer_.expect_write_aliased_string_(value);
    end

    function write_aliased_enum_(obj, value)
      obj.writer_.write_aliased_enum(value);
      obj.mock_writer_.expect_write_aliased_enum_(value);
    end

    function write_aliased_open_generic_(obj, value)
      obj.writer_.write_aliased_open_generic(value);
      obj.mock_writer_.expect_write_aliased_open_generic_(value);
    end

    function write_aliased_closed_generic_(obj, value)
      obj.writer_.write_aliased_closed_generic(value);
      obj.mock_writer_.expect_write_aliased_closed_generic_(value);
    end

    function write_aliased_optional_(obj, value)
      obj.writer_.write_aliased_optional(value);
      obj.mock_writer_.expect_write_aliased_optional_(value);
    end

    function write_aliased_generic_optional_(obj, value)
      obj.writer_.write_aliased_generic_optional(value);
      obj.mock_writer_.expect_write_aliased_generic_optional_(value);
    end

    function write_aliased_generic_union_2_(obj, value)
      obj.writer_.write_aliased_generic_union_2(value);
      obj.mock_writer_.expect_write_aliased_generic_union_2_(value);
    end

    function write_aliased_generic_vector_(obj, value)
      obj.writer_.write_aliased_generic_vector(value);
      obj.mock_writer_.expect_write_aliased_generic_vector_(value);
    end

    function write_aliased_generic_fixed_vector_(obj, value)
      obj.writer_.write_aliased_generic_fixed_vector(value);
      obj.mock_writer_.expect_write_aliased_generic_fixed_vector_(value);
    end

    function write_stream_of_aliased_generic_union_2_(obj, value)
      obj.writer_.write_stream_of_aliased_generic_union_2(value);
      obj.mock_writer_.expect_write_stream_of_aliased_generic_union_2_(value);
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
