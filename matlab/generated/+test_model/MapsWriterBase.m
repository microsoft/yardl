% Abstract writer for protocol Maps
classdef (Abstract) MapsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = MapsWriterBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 8
        expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_string_to_int(obj, value)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_string_to_int_(value);
      obj.state_ = 2;
    end

    % Ordinal 1
    function write_int_to_string(obj, value)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      obj.write_int_to_string_(value);
      obj.state_ = 4;
    end

    % Ordinal 2
    function write_string_to_union(obj, value)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      obj.write_string_to_union_(value);
      obj.state_ = 6;
    end

    % Ordinal 3
    function write_aliased_generic(obj, value)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      obj.write_aliased_generic_(value);
      obj.state_ = 8;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"Maps","sequence":[{"name":"stringToInt","type":{"map":{"keys":"string","values":"int32"}}},{"name":"intToString","type":{"map":{"keys":"int32","values":"string"}}},{"name":"stringToUnion","type":{"map":{"keys":"string","values":[{"tag":"string","type":"string"},{"tag":"int32","type":"int32"}]}}},{"name":"aliasedGeneric","type":{"name":"BasicTypes.AliasedMap","typeArguments":["string","int32"]}}]},"types":[{"name":"AliasedMap","typeParameters":["K","V"],"type":{"map":{"keys":"K","values":"V"}}}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_string_to_int_(obj, value)
    write_int_to_string_(obj, value)
    write_string_to_union_(obj, value)
    write_aliased_generic_(obj, value)

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
        name = 'write_string_to_int';
      elseif state == 2
        name = 'write_int_to_string';
      elseif state == 4
        name = 'write_string_to_union';
      elseif state == 6
        name = 'write_aliased_generic';
      else
        name = '<unknown>';
      end
    end
  end
end
