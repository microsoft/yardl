% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef SubarraysInRecordsReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function self = SubarraysInRecordsReaderBase()
      self.state_ = 0;
    end

    function close(self)
      self.close_();
      if self.state_ ~= 2
        expected_method = self.state_to_method_name_(self.state_);
        throw(yardl.ProtocolError("Protocol reader closed before all data was consumed. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function value = read_with_fixed_subarrays(self)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      value = self.read_with_fixed_subarrays_();
      self.state_ = 1;
    end

    % Ordinal 1
    function value = read_with_vlen_subarrays(self)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      value = self.read_with_vlen_subarrays_();
      self.state_ = 2;
    end

    function copy_to(self, writer)
      writer.write_with_fixed_subarrays(self.read_with_fixed_subarrays());
      writer.write_with_vlen_subarrays(self.read_with_vlen_subarrays());
    end
  end

  methods (Static)
    function res = schema()
      res = test_model.SubarraysInRecordsWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    read_with_fixed_subarrays_(self)
    read_with_vlen_subarrays_(self)

    close_(self)
  end

  methods (Access=private)
    function raise_unexpected_state_(self, actual)
      actual_method = self.state_to_method_name_(actual);
      expected_method = self.state_to_method_name_(self.state_);
      throw(yardl.ProtocolError("Expected call to '%s' but received call to '%s'.", expected_method, actual_method));
    end

    function name = state_to_method_name_(self, state)
      if state == 0
        name = "read_with_fixed_subarrays";
      elseif state == 1
        name = "read_with_vlen_subarrays";
      else
        name = "<unknown>";
      end
    end
  end
end
