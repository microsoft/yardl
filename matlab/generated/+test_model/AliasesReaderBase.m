classdef AliasesReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = AliasesReaderBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 20
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
    function value = read_aliased_string(obj)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      value = obj.read_aliased_string_();
      obj.state_ = 2;
    end

    % Ordinal 1
    function value = read_aliased_enum(obj)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      value = obj.read_aliased_enum_();
      obj.state_ = 4;
    end

    % Ordinal 2
    function value = read_aliased_open_generic(obj)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      value = obj.read_aliased_open_generic_();
      obj.state_ = 6;
    end

    % Ordinal 3
    function value = read_aliased_closed_generic(obj)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      value = obj.read_aliased_closed_generic_();
      obj.state_ = 8;
    end

    % Ordinal 4
    function value = read_aliased_optional(obj)
      if obj.state_ ~= 8
        obj.raise_unexpected_state_(8);
      end

      value = obj.read_aliased_optional_();
      obj.state_ = 10;
    end

    % Ordinal 5
    function value = read_aliased_generic_optional(obj)
      if obj.state_ ~= 10
        obj.raise_unexpected_state_(10);
      end

      value = obj.read_aliased_generic_optional_();
      obj.state_ = 12;
    end

    % Ordinal 6
    function value = read_aliased_generic_union_2(obj)
      if obj.state_ ~= 12
        obj.raise_unexpected_state_(12);
      end

      value = obj.read_aliased_generic_union_2_();
      obj.state_ = 14;
    end

    % Ordinal 7
    function value = read_aliased_generic_vector(obj)
      if obj.state_ ~= 14
        obj.raise_unexpected_state_(14);
      end

      value = obj.read_aliased_generic_vector_();
      obj.state_ = 16;
    end

    % Ordinal 8
    function value = read_aliased_generic_fixed_vector(obj)
      if obj.state_ ~= 16
        obj.raise_unexpected_state_(16);
      end

      value = obj.read_aliased_generic_fixed_vector_();
      obj.state_ = 18;
    end

    % Ordinal 9
    function value = read_stream_of_aliased_generic_union_2(obj)
      if obj.state_ ~= 18
        obj.raise_unexpected_state_(18);
      end

      value = obj.read_stream_of_aliased_generic_union_2_();
      obj.state_ = 20;
    end

    function copy_to(obj, writer)
      writer.write_aliased_string(obj.read_aliased_string());
      writer.write_aliased_enum(obj.read_aliased_enum());
      writer.write_aliased_open_generic(obj.read_aliased_open_generic());
      writer.write_aliased_closed_generic(obj.read_aliased_closed_generic());
      writer.write_aliased_optional(obj.read_aliased_optional());
      writer.write_aliased_generic_optional(obj.read_aliased_generic_optional());
      writer.write_aliased_generic_union_2(obj.read_aliased_generic_union_2());
      writer.write_aliased_generic_vector(obj.read_aliased_generic_vector());
      writer.write_aliased_generic_fixed_vector(obj.read_aliased_generic_fixed_vector());
      writer.write_stream_of_aliased_generic_union_2(obj.read_stream_of_aliased_generic_union_2());
    end
  end

  methods (Static)
    function res = schema()
      res = test_model.AliasesWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    read_aliased_string_(obj, value)
    read_aliased_enum_(obj, value)
    read_aliased_open_generic_(obj, value)
    read_aliased_closed_generic_(obj, value)
    read_aliased_optional_(obj, value)
    read_aliased_generic_optional_(obj, value)
    read_aliased_generic_union_2_(obj, value)
    read_aliased_generic_vector_(obj, value)
    read_aliased_generic_fixed_vector_(obj, value)
    read_stream_of_aliased_generic_union_2_(obj, value)

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
        name = 'read_aliased_string';
      elseif state == 2
        name = 'read_aliased_enum';
      elseif state == 4
        name = 'read_aliased_open_generic';
      elseif state == 6
        name = 'read_aliased_closed_generic';
      elseif state == 8
        name = 'read_aliased_optional';
      elseif state == 10
        name = 'read_aliased_generic_optional';
      elseif state == 12
        name = 'read_aliased_generic_union_2';
      elseif state == 14
        name = 'read_aliased_generic_vector';
      elseif state == 16
        name = 'read_aliased_generic_fixed_vector';
      elseif state == 18
        name = 'read_stream_of_aliased_generic_union_2';
      else
        name = '<unknown>';
      end
    end
  end
end
