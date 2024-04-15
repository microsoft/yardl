% This file was generated by the "yardl" tool. DO NOT EDIT.

% Abstract writer for protocol ProtocolWithKeywordSteps
classdef (Abstract) ProtocolWithKeywordStepsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function self = ProtocolWithKeywordStepsWriterBase()
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
    function write_int(self, value)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.write_int_(value);
    end

    function end_int(self)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.end_stream_();
      self.state_ = 1;
    end

    % Ordinal 1
    function write_float(self, value)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      self.write_float_(value);
      self.state_ = 2;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"ProtocolWithKeywordSteps","sequence":[{"name":"int","type":{"stream":{"items":"TestModel.RecordWithKeywordFields"}}},{"name":"float","type":"TestModel.EnumWithKeywordSymbols"}]},"types":[{"name":"ArrayWithKeywordDimensionNames","type":{"array":{"items":"int32","dimensions":[{"name":"while"},{"name":"do"}]}}},{"name":"EnumWithKeywordSymbols","values":[{"symbol":"try","value":2},{"symbol":"catch","value":1}]},{"name":"RecordWithKeywordFields","fields":[{"name":"int","type":"string"},{"name":"sizeof","type":"TestModel.ArrayWithKeywordDimensionNames"},{"name":"if","type":"TestModel.EnumWithKeywordSymbols"}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_int_(self, value)
    write_float_(self, value)

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
        name = "write_int or end_int";
      elseif state == 1
        name = "write_float";
      else
        name = '<unknown>';
      end
    end
  end
end
