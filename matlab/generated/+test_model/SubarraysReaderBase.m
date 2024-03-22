classdef SubarraysReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = SubarraysReaderBase()
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
    function value = read_dynamic_with_fixed_int_subarray(obj)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      value = obj.read_dynamic_with_fixed_int_subarray_();
      obj.state_ = 2;
    end

    % Ordinal 1
    function value = read_dynamic_with_fixed_float_subarray(obj)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      value = obj.read_dynamic_with_fixed_float_subarray_();
      obj.state_ = 4;
    end

    % Ordinal 2
    function value = read_known_dim_count_with_fixed_int_subarray(obj)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      value = obj.read_known_dim_count_with_fixed_int_subarray_();
      obj.state_ = 6;
    end

    % Ordinal 3
    function value = read_known_dim_count_with_fixed_float_subarray(obj)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      value = obj.read_known_dim_count_with_fixed_float_subarray_();
      obj.state_ = 8;
    end

    % Ordinal 4
    function value = read_fixed_with_fixed_int_subarray(obj)
      if obj.state_ ~= 8
        obj.raise_unexpected_state_(8);
      end

      value = obj.read_fixed_with_fixed_int_subarray_();
      obj.state_ = 10;
    end

    % Ordinal 5
    function value = read_fixed_with_fixed_float_subarray(obj)
      if obj.state_ ~= 10
        obj.raise_unexpected_state_(10);
      end

      value = obj.read_fixed_with_fixed_float_subarray_();
      obj.state_ = 12;
    end

    % Ordinal 6
    function value = read_nested_subarray(obj)
      if obj.state_ ~= 12
        obj.raise_unexpected_state_(12);
      end

      value = obj.read_nested_subarray_();
      obj.state_ = 14;
    end

    % Ordinal 7
    function value = read_dynamic_with_fixed_vector_subarray(obj)
      if obj.state_ ~= 14
        obj.raise_unexpected_state_(14);
      end

      value = obj.read_dynamic_with_fixed_vector_subarray_();
      obj.state_ = 16;
    end

    % Ordinal 8
    function value = read_generic_subarray(obj)
      if obj.state_ ~= 16
        obj.raise_unexpected_state_(16);
      end

      value = obj.read_generic_subarray_();
      obj.state_ = 18;
    end

    function copy_to(obj, writer)
      writer.write_dynamic_with_fixed_int_subarray(obj.read_dynamic_with_fixed_int_subarray());
      writer.write_dynamic_with_fixed_float_subarray(obj.read_dynamic_with_fixed_float_subarray());
      writer.write_known_dim_count_with_fixed_int_subarray(obj.read_known_dim_count_with_fixed_int_subarray());
      writer.write_known_dim_count_with_fixed_float_subarray(obj.read_known_dim_count_with_fixed_float_subarray());
      writer.write_fixed_with_fixed_int_subarray(obj.read_fixed_with_fixed_int_subarray());
      writer.write_fixed_with_fixed_float_subarray(obj.read_fixed_with_fixed_float_subarray());
      writer.write_nested_subarray(obj.read_nested_subarray());
      writer.write_dynamic_with_fixed_vector_subarray(obj.read_dynamic_with_fixed_vector_subarray());
      writer.write_generic_subarray(obj.read_generic_subarray());
    end
  end

  methods (Static)
    function res = schema()
      res = test_model.SubarraysWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    read_dynamic_with_fixed_int_subarray_(obj, value)
    read_dynamic_with_fixed_float_subarray_(obj, value)
    read_known_dim_count_with_fixed_int_subarray_(obj, value)
    read_known_dim_count_with_fixed_float_subarray_(obj, value)
    read_fixed_with_fixed_int_subarray_(obj, value)
    read_fixed_with_fixed_float_subarray_(obj, value)
    read_nested_subarray_(obj, value)
    read_dynamic_with_fixed_vector_subarray_(obj, value)
    read_generic_subarray_(obj, value)

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
        name = 'read_dynamic_with_fixed_int_subarray';
      elseif state == 2
        name = 'read_dynamic_with_fixed_float_subarray';
      elseif state == 4
        name = 'read_known_dim_count_with_fixed_int_subarray';
      elseif state == 6
        name = 'read_known_dim_count_with_fixed_float_subarray';
      elseif state == 8
        name = 'read_fixed_with_fixed_int_subarray';
      elseif state == 10
        name = 'read_fixed_with_fixed_float_subarray';
      elseif state == 12
        name = 'read_nested_subarray';
      elseif state == 14
        name = 'read_dynamic_with_fixed_vector_subarray';
      elseif state == 16
        name = 'read_generic_subarray';
      else
        name = '<unknown>';
      end
    end
  end
end