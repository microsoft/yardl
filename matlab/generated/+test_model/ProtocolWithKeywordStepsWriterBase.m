% Abstract writer for protocol ProtocolWithKeywordSteps
classdef (Abstract) ProtocolWithKeywordStepsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = ProtocolWithKeywordStepsWriterBase()
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
    function write_int(obj, value)
      if bitand(int32(obj.state_), bitcmp(1, 'int8')) ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_int_(value);
      obj.state_ = 1;
    end

    % Ordinal 1
    function write_float(obj, value)
      if obj.state_ == 1
        obj.end_stream_();
        obj.state_ = 2;
      elseif obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      obj.write_float_(value);
      obj.state_ = 4;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"ProtocolWithKeywordSteps","sequence":[{"name":"int","type":{"stream":{"items":"TestModel.RecordWithKeywordFields"}}},{"name":"float","type":"TestModel.EnumWithKeywordSymbols"}]},"types":[{"name":"ArrayWithKeywordDimensionNames","type":{"array":{"items":"int32","dimensions":[{"name":"while"},{"name":"do"}]}}},{"name":"EnumWithKeywordSymbols","values":[{"symbol":"try","value":2},{"symbol":"catch","value":1}]},{"name":"RecordWithKeywordFields","fields":[{"name":"int","type":"string"},{"name":"sizeof","type":"TestModel.ArrayWithKeywordDimensionNames"},{"name":"if","type":"TestModel.EnumWithKeywordSymbols"}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_int_(obj, value)
    write_float_(obj, value)

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
        name = 'write_int';
      elseif state == 2
        name = 'write_float';
      else
        name = '<unknown>';
      end
    end
  end
end
