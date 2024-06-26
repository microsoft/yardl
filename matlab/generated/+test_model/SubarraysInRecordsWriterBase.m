% This file was generated by the "yardl" tool. DO NOT EDIT.

% Abstract writer for protocol SubarraysInRecords
classdef (Abstract) SubarraysInRecordsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function self = SubarraysInRecordsWriterBase()
      self.state_ = 0;
    end

    function close(self)
      self.close_();
      if self.state_ ~= 2
        expected_method = self.state_to_method_name_(self.state_);
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_with_fixed_subarrays(self, value)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.write_with_fixed_subarrays_(value);
      self.state_ = 1;
    end

    % Ordinal 1
    function write_with_vlen_subarrays(self, value)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      self.write_with_vlen_subarrays_(value);
      self.state_ = 2;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"SubarraysInRecords","sequence":[{"name":"withFixedSubarrays","type":{"array":{"items":"TestModel.RecordWithFixedCollections"}}},{"name":"withVlenSubarrays","type":{"array":{"items":"TestModel.RecordWithVlenCollections"}}}]},"types":[{"name":"RecordWithFixedCollections","fields":[{"name":"fixedVector","type":{"vector":{"items":"int32","length":3}}},{"name":"fixedArray","type":{"array":{"items":"int32","dimensions":[{"length":2},{"length":3}]}}}]},{"name":"RecordWithVlenCollections","fields":[{"name":"vector","type":{"vector":{"items":"int32"}}},{"name":"array","type":{"array":{"items":"int32","dimensions":2}}}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_with_fixed_subarrays_(self, value)
    write_with_vlen_subarrays_(self, value)

    end_stream_(self)
    close_(self)
  end

  methods (Access=private)
    function raise_unexpected_state_(self, actual)
      expected_method = self.state_to_method_name_(self.state_);
      actual_method = self.state_to_method_name_(actual);
      throw(yardl.ProtocolError("Expected call to '%s' but received call to '%s'", expected_method, actual_method));
    end

    function name = state_to_method_name_(self, state)
      if state == 0
        name = "write_with_fixed_subarrays";
      elseif state == 1
        name = "write_with_vlen_subarrays";
      else
        name = '<unknown>';
      end
    end
  end
end
