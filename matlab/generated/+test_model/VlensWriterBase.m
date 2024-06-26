% This file was generated by the "yardl" tool. DO NOT EDIT.

% Abstract writer for protocol Vlens
classdef (Abstract) VlensWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function self = VlensWriterBase()
      self.state_ = 0;
    end

    function close(self)
      self.close_();
      if self.state_ ~= 4
        expected_method = self.state_to_method_name_(self.state_);
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_int_vector(self, value)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.write_int_vector_(value);
      self.state_ = 1;
    end

    % Ordinal 1
    function write_complex_vector(self, value)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      self.write_complex_vector_(value);
      self.state_ = 2;
    end

    % Ordinal 2
    function write_record_with_vlens(self, value)
      if self.state_ ~= 2
        self.raise_unexpected_state_(2);
      end

      self.write_record_with_vlens_(value);
      self.state_ = 3;
    end

    % Ordinal 3
    function write_vlen_of_record_with_vlens(self, value)
      if self.state_ ~= 3
        self.raise_unexpected_state_(3);
      end

      self.write_vlen_of_record_with_vlens_(value);
      self.state_ = 4;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"Vlens","sequence":[{"name":"intVector","type":{"vector":{"items":"int32"}}},{"name":"complexVector","type":{"vector":{"items":"complexfloat32"}}},{"name":"recordWithVlens","type":"TestModel.RecordWithVlens"},{"name":"vlenOfRecordWithVlens","type":{"vector":{"items":"TestModel.RecordWithVlens"}}}]},"types":[{"name":"RecordWithVlens","fields":[{"name":"a","type":{"vector":{"items":"TestModel.SimpleRecord"}}},{"name":"b","type":"int32"},{"name":"c","type":"int32"}]},{"name":"SimpleRecord","fields":[{"name":"x","type":"int32"},{"name":"y","type":"int32"},{"name":"z","type":"int32"}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_int_vector_(self, value)
    write_complex_vector_(self, value)
    write_record_with_vlens_(self, value)
    write_vlen_of_record_with_vlens_(self, value)

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
        name = "write_int_vector";
      elseif state == 1
        name = "write_complex_vector";
      elseif state == 2
        name = "write_record_with_vlens";
      elseif state == 3
        name = "write_vlen_of_record_with_vlens";
      else
        name = '<unknown>';
      end
    end
  end
end
