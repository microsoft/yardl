classdef MapsReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = MapsReaderBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 8
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
    function value = read_string_to_int(obj)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      value = obj.read_string_to_int_();
      obj.state_ = 2;
    end

    % Ordinal 1
    function value = read_int_to_string(obj)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      value = obj.read_int_to_string_();
      obj.state_ = 4;
    end

    % Ordinal 2
    function value = read_string_to_union(obj)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      value = obj.read_string_to_union_();
      obj.state_ = 6;
    end

    % Ordinal 3
    function value = read_aliased_generic(obj)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      value = obj.read_aliased_generic_();
      obj.state_ = 8;
    end

    function copy_to(obj, writer)
      writer.write_string_to_int(obj.read_string_to_int());
      writer.write_int_to_string(obj.read_int_to_string());
      writer.write_string_to_union(obj.read_string_to_union());
      writer.write_aliased_generic(obj.read_aliased_generic());
    end
  end

  methods (Static)
    function res = schema()
      res = test_model.MapsWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    read_string_to_int_(obj, value)
    read_int_to_string_(obj, value)
    read_string_to_union_(obj, value)
    read_aliased_generic_(obj, value)

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
        name = 'read_string_to_int';
      elseif state == 2
        name = 'read_int_to_string';
      elseif state == 4
        name = 'read_string_to_union';
      elseif state == 6
        name = 'read_aliased_generic';
      else
        name = '<unknown>';
      end
    end
  end
end
