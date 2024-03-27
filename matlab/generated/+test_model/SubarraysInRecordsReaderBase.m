% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef SubarraysInRecordsReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = SubarraysInRecordsReaderBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 4
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
    function value = read_with_fixed_subarrays(obj)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      value = obj.read_with_fixed_subarrays_();
      obj.state_ = 2;
    end

    % Ordinal 1
    function value = read_with_vlen_subarrays(obj)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      value = obj.read_with_vlen_subarrays_();
      obj.state_ = 4;
    end

    function copy_to(obj, writer)
      writer.write_with_fixed_subarrays(obj.read_with_fixed_subarrays());
      writer.write_with_vlen_subarrays(obj.read_with_vlen_subarrays());
    end
  end

  methods (Static)
    function res = schema()
      res = test_model.SubarraysInRecordsWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    read_with_fixed_subarrays_(obj, value)
    read_with_vlen_subarrays_(obj, value)

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
        name = 'read_with_fixed_subarrays';
      elseif state == 2
        name = 'read_with_vlen_subarrays';
      else
        name = '<unknown>';
      end
    end
  end
end
