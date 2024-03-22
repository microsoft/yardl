% Abstract writer for protocol SubarraysInRecords
classdef (Abstract) SubarraysInRecordsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = SubarraysInRecordsWriterBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 4
        expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_with_fixed_subarrays(obj, value)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_with_fixed_subarrays_(value);
      obj.state_ = 2;
    end

    % Ordinal 1
    function write_with_vlen_subarrays(obj, value)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      obj.write_with_vlen_subarrays_(value);
      obj.state_ = 4;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"SubarraysInRecords","sequence":[{"name":"withFixedSubarrays","type":{"array":{"items":"TestModel.RecordWithFixedCollections"}}},{"name":"withVlenSubarrays","type":{"array":{"items":"TestModel.RecordWithVlenCollections"}}}]},"types":[{"name":"RecordWithFixedCollections","fields":[{"name":"fixedVector","type":{"vector":{"items":"int32","length":3}}},{"name":"fixedArray","type":{"array":{"items":"int32","dimensions":[{"length":2},{"length":3}]}}}]},{"name":"RecordWithVlenCollections","fields":[{"name":"vector","type":{"vector":{"items":"int32"}}},{"name":"array","type":{"array":{"items":"int32","dimensions":2}}}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_with_fixed_subarrays_(obj, value)
    write_with_vlen_subarrays_(obj, value)

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
        name = 'write_with_fixed_subarrays';
      elseif state == 2
        name = 'write_with_vlen_subarrays';
      else
        name = '<unknown>';
      end
    end
  end
end
