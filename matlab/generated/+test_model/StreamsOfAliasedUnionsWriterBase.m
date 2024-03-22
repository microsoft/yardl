% Abstract writer for protocol StreamsOfAliasedUnions
classdef (Abstract) StreamsOfAliasedUnionsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = StreamsOfAliasedUnionsWriterBase()
      obj.state_ = 0;
    end

    function close(obj)
      if obj.state_ == 3
        obj.end_stream_();
        obj.close_();
        return
      end
      obj.close_();
      if obj.state_ ~= 4
        expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_int_or_simple_record(obj, value)
      if bitand(int32(obj.state_), bitcmp(1, 'int8')) ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_int_or_simple_record_(value);
      obj.state_ = 1;
    end

    % Ordinal 1
    function write_nullable_int_or_simple_record(obj, value)
      if obj.state_ == 1
        obj.end_stream_();
        obj.state_ = 2;
      elseif bitand(int32(obj.state_), bitcmp(1, 'int8')) ~= 2
        obj.raise_unexpected_state_(2);
      end

      obj.write_nullable_int_or_simple_record_(value);
      obj.state_ = 3;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"StreamsOfAliasedUnions","sequence":[{"name":"intOrSimpleRecord","type":{"stream":{"items":"TestModel.AliasedIntOrSimpleRecord"}}},{"name":"nullableIntOrSimpleRecord","type":{"stream":{"items":"TestModel.AliasedNullableIntSimpleRecord"}}}]},"types":[{"name":"AliasedIntOrSimpleRecord","type":[{"tag":"int32","type":"int32"},{"tag":"SimpleRecord","type":"TestModel.SimpleRecord"}]},{"name":"AliasedNullableIntSimpleRecord","type":[null,{"tag":"int32","type":"int32"},{"tag":"SimpleRecord","type":"TestModel.SimpleRecord"}]},{"name":"SimpleRecord","fields":[{"name":"x","type":"int32"},{"name":"y","type":"int32"},{"name":"z","type":"int32"}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_int_or_simple_record_(obj, value)
    write_nullable_int_or_simple_record_(obj, value)

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
        name = 'write_int_or_simple_record';
      elseif state == 2
        name = 'write_nullable_int_or_simple_record';
      else
        name = '<unknown>';
      end
    end
  end
end
