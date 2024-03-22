% Abstract writer for protocol StateTest
classdef (Abstract) StateTestWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = StateTestWriterBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 6
        expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_an_int(obj, value)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_an_int_(value);
      obj.state_ = 2;
    end

    % Ordinal 1
    function write_a_stream(obj, value)
      if bitand(int32(obj.state_), bitcmp(1, 'int8')) ~= 2
        obj.raise_unexpected_state_(2);
      end

      obj.write_a_stream_(value);
      obj.state_ = 3;
    end

    % Ordinal 2
    function write_another_int(obj, value)
      if obj.state_ == 3
        obj.end_stream_();
        obj.state_ = 4;
      elseif obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      obj.write_another_int_(value);
      obj.state_ = 6;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"StateTest","sequence":[{"name":"anInt","type":"int32"},{"name":"aStream","type":{"stream":{"items":"int32"}}},{"name":"anotherInt","type":"int32"}]},"types":null}');
    end
  end

  methods (Abstract, Access=protected)
    write_an_int_(obj, value)
    write_a_stream_(obj, value)
    write_another_int_(obj, value)

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
        name = 'write_an_int';
      elseif state == 2
        name = 'write_a_stream';
      elseif state == 4
        name = 'write_another_int';
      else
        name = '<unknown>';
      end
    end
  end
end
