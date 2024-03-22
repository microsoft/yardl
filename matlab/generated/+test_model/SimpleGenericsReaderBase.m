classdef SimpleGenericsReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = SimpleGenericsReaderBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 18
        if mod(obj.state_, 2) == 1
          previous_method = obj.state_to_method_name_(obj.state_ - 1);
          throw(yardl.ProtocolError("Protocol reader closed before all data was consumed. The iterable returned by '%s' was not fully consumed.", previous_method));
        else
          expected_method = obj.state_to_method_name_(obj.state_);
          throw(yardl.ProtocolError("Protocol reader closed before all data was consumed. Expected call to '%s'.", expected_method));
        end
      end
    end

    % Ordinal 0
    function value = read_float_image(obj)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      value = obj.read_float_image_();
      obj.state_ = 2;
    end

    % Ordinal 1
    function value = read_int_image(obj)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      value = obj.read_int_image_();
      obj.state_ = 4;
    end

    % Ordinal 2
    function value = read_int_image_alternate_syntax(obj)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      value = obj.read_int_image_alternate_syntax_();
      obj.state_ = 6;
    end

    % Ordinal 3
    function value = read_string_image(obj)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      value = obj.read_string_image_();
      obj.state_ = 8;
    end

    % Ordinal 4
    function value = read_int_float_tuple(obj)
      if obj.state_ ~= 8
        obj.raise_unexpected_state_(8);
      end

      value = obj.read_int_float_tuple_();
      obj.state_ = 10;
    end

    % Ordinal 5
    function value = read_float_float_tuple(obj)
      if obj.state_ ~= 10
        obj.raise_unexpected_state_(10);
      end

      value = obj.read_float_float_tuple_();
      obj.state_ = 12;
    end

    % Ordinal 6
    function value = read_int_float_tuple_alternate_syntax(obj)
      if obj.state_ ~= 12
        obj.raise_unexpected_state_(12);
      end

      value = obj.read_int_float_tuple_alternate_syntax_();
      obj.state_ = 14;
    end

    % Ordinal 7
    function value = read_int_string_tuple(obj)
      if obj.state_ ~= 14
        obj.raise_unexpected_state_(14);
      end

      value = obj.read_int_string_tuple_();
      obj.state_ = 16;
    end

    % Ordinal 8
    function value = read_stream_of_type_variants(obj)
      if obj.state_ ~= 16
        obj.raise_unexpected_state_(16);
      end

      value = obj.read_stream_of_type_variants_();
      obj.state_ = 18;
    end

    function copy_to(obj, writer)
      writer.write_float_image(obj.read_float_image());
      writer.write_int_image(obj.read_int_image());
      writer.write_int_image_alternate_syntax(obj.read_int_image_alternate_syntax());
      writer.write_string_image(obj.read_string_image());
      writer.write_int_float_tuple(obj.read_int_float_tuple());
      writer.write_float_float_tuple(obj.read_float_float_tuple());
      writer.write_int_float_tuple_alternate_syntax(obj.read_int_float_tuple_alternate_syntax());
      writer.write_int_string_tuple(obj.read_int_string_tuple());
      writer.write_stream_of_type_variants(obj.read_stream_of_type_variants());
    end
  end

  methods (Static)
    function res = schema()
      res = test_model.SimpleGenericsWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    read_float_image_(obj, value)
    read_int_image_(obj, value)
    read_int_image_alternate_syntax_(obj, value)
    read_string_image_(obj, value)
    read_int_float_tuple_(obj, value)
    read_float_float_tuple_(obj, value)
    read_int_float_tuple_alternate_syntax_(obj, value)
    read_int_string_tuple_(obj, value)
    read_stream_of_type_variants_(obj, value)

    close_(obj)
  end

  methods (Access=private)
    function raise_unexpected_state_(obj, actual)
      actual_method = obj.state_to_method_name_(actual);
      expected_method = obj.state_to_method_name_(obj.state_);
      throw(yardl.ProtocolError("Expected call to '%s' but received call to '%s'.", expected_method, actual_method));
    end

    function name = state_to_method_name_(obj, state)
      if state == 0
        name = 'read_float_image';
      elseif state == 2
        name = 'read_int_image';
      elseif state == 4
        name = 'read_int_image_alternate_syntax';
      elseif state == 6
        name = 'read_string_image';
      elseif state == 8
        name = 'read_int_float_tuple';
      elseif state == 10
        name = 'read_float_float_tuple';
      elseif state == 12
        name = 'read_int_float_tuple_alternate_syntax';
      elseif state == 14
        name = 'read_int_string_tuple';
      elseif state == 16
        name = 'read_stream_of_type_variants';
      else
        name = '<unknown>';
      end
    end
  end
end