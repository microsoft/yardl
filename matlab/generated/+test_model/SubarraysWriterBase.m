% Abstract writer for protocol Subarrays
classdef (Abstract) SubarraysWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = SubarraysWriterBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 18
        expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_dynamic_with_fixed_int_subarray(obj, value)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_dynamic_with_fixed_int_subarray_(value);
      obj.state_ = 2;
    end

    % Ordinal 1
    function write_dynamic_with_fixed_float_subarray(obj, value)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      obj.write_dynamic_with_fixed_float_subarray_(value);
      obj.state_ = 4;
    end

    % Ordinal 2
    function write_known_dim_count_with_fixed_int_subarray(obj, value)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      obj.write_known_dim_count_with_fixed_int_subarray_(value);
      obj.state_ = 6;
    end

    % Ordinal 3
    function write_known_dim_count_with_fixed_float_subarray(obj, value)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      obj.write_known_dim_count_with_fixed_float_subarray_(value);
      obj.state_ = 8;
    end

    % Ordinal 4
    function write_fixed_with_fixed_int_subarray(obj, value)
      if obj.state_ ~= 8
        obj.raise_unexpected_state_(8);
      end

      obj.write_fixed_with_fixed_int_subarray_(value);
      obj.state_ = 10;
    end

    % Ordinal 5
    function write_fixed_with_fixed_float_subarray(obj, value)
      if obj.state_ ~= 10
        obj.raise_unexpected_state_(10);
      end

      obj.write_fixed_with_fixed_float_subarray_(value);
      obj.state_ = 12;
    end

    % Ordinal 6
    function write_nested_subarray(obj, value)
      if obj.state_ ~= 12
        obj.raise_unexpected_state_(12);
      end

      obj.write_nested_subarray_(value);
      obj.state_ = 14;
    end

    % Ordinal 7
    function write_dynamic_with_fixed_vector_subarray(obj, value)
      if obj.state_ ~= 14
        obj.raise_unexpected_state_(14);
      end

      obj.write_dynamic_with_fixed_vector_subarray_(value);
      obj.state_ = 16;
    end

    % Ordinal 8
    function write_generic_subarray(obj, value)
      if obj.state_ ~= 16
        obj.raise_unexpected_state_(16);
      end

      obj.write_generic_subarray_(value);
      obj.state_ = 18;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"Subarrays","sequence":[{"name":"dynamicWithFixedIntSubarray","type":{"array":{"items":{"array":{"items":"int32","dimensions":[{"length":3}]}}}}},{"name":"dynamicWithFixedFloatSubarray","type":{"array":{"items":{"array":{"items":"float32","dimensions":[{"length":3}]}}}}},{"name":"knownDimCountWithFixedIntSubarray","type":{"array":{"items":{"array":{"items":"int32","dimensions":[{"length":3}]}},"dimensions":1}}},{"name":"knownDimCountWithFixedFloatSubarray","type":{"array":{"items":{"array":{"items":"float32","dimensions":[{"length":3}]}},"dimensions":1}}},{"name":"fixedWithFixedIntSubarray","type":{"array":{"items":{"array":{"items":"int32","dimensions":[{"length":3}]}},"dimensions":[{"length":2}]}}},{"name":"fixedWithFixedFloatSubarray","type":{"array":{"items":{"array":{"items":"float32","dimensions":[{"length":3}]}},"dimensions":[{"length":2}]}}},{"name":"nestedSubarray","type":{"array":{"items":{"array":{"items":{"array":{"items":"int32","dimensions":[{"length":3}]}},"dimensions":[{"length":2}]}}}}},{"name":"dynamicWithFixedVectorSubarray","type":{"array":{"items":{"vector":{"items":"int32","length":3}}}}},{"name":"genericSubarray","type":{"name":"TestModel.Image","typeArguments":[{"array":{"items":"int32","dimensions":[{"length":3}]}}]}}]},"types":[{"name":"Image","typeParameters":["T"],"type":{"array":{"items":"T","dimensions":[{"name":"x"},{"name":"y"}]}}},{"name":"Image","typeParameters":["T"],"type":{"name":"Image.Image","typeArguments":["T"]}}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_dynamic_with_fixed_int_subarray_(obj, value)
    write_dynamic_with_fixed_float_subarray_(obj, value)
    write_known_dim_count_with_fixed_int_subarray_(obj, value)
    write_known_dim_count_with_fixed_float_subarray_(obj, value)
    write_fixed_with_fixed_int_subarray_(obj, value)
    write_fixed_with_fixed_float_subarray_(obj, value)
    write_nested_subarray_(obj, value)
    write_dynamic_with_fixed_vector_subarray_(obj, value)
    write_generic_subarray_(obj, value)

    end_stream_(obj)
    close_(obj)
  end

  methods (Access=private)
    function raise_unexpected_state_(obj, actual)
      expected_method = obj.state_to_method_name_(obj.state_);
      actual_method = obj.state_to_method_name_(actual);
      throw(yardl.ProtocolError("Expected call to '%s' but received call to '%s'", expected_method, actual_method));
    end

    function name = state_to_method_name_(obj, state)
      if state == 0
        name = 'write_dynamic_with_fixed_int_subarray';
      elseif state == 2
        name = 'write_dynamic_with_fixed_float_subarray';
      elseif state == 4
        name = 'write_known_dim_count_with_fixed_int_subarray';
      elseif state == 6
        name = 'write_known_dim_count_with_fixed_float_subarray';
      elseif state == 8
        name = 'write_fixed_with_fixed_int_subarray';
      elseif state == 10
        name = 'write_fixed_with_fixed_float_subarray';
      elseif state == 12
        name = 'write_nested_subarray';
      elseif state == 14
        name = 'write_dynamic_with_fixed_vector_subarray';
      elseif state == 16
        name = 'write_generic_subarray';
      else
        name = '<unknown>';
      end
    end
  end
end
