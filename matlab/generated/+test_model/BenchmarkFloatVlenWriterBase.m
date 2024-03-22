% Abstract writer for protocol BenchmarkFloatVlen
classdef (Abstract) BenchmarkFloatVlenWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = BenchmarkFloatVlenWriterBase()
      obj.state_ = 0;
    end

    function close(obj)
      if obj.state_ == 1
        obj.end_stream_();
        obj.close_();
        return
      end
      obj.close_();
      if obj.state_ ~= 2
        expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_float_array(obj, value)
      if bitand(int32(obj.state_), bitcmp(1, 'int8')) ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_float_array_(value);
      obj.state_ = 1;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"BenchmarkFloatVlen","sequence":[{"name":"floatArray","type":{"stream":{"items":{"array":{"items":"float32","dimensions":2}}}}}]},"types":null}');
    end
  end

  methods (Abstract, Access=protected)
    write_float_array_(obj, value)

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
        name = 'write_float_array';
      else
        name = '<unknown>';
      end
    end
  end
end