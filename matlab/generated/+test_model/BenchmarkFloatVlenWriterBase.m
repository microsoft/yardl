% This file was generated by the "yardl" tool. DO NOT EDIT.

% Abstract writer for protocol BenchmarkFloatVlen
classdef (Abstract) BenchmarkFloatVlenWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function self = BenchmarkFloatVlenWriterBase()
      self.state_ = 0;
    end

    function close(self)
      self.close_();
      if self.state_ ~= 1
        expected_method = self.state_to_method_name_(self.state_);
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_float_array(self, value)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.write_float_array_(value);
    end

    function end_float_array(self)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.end_stream_();
      self.state_ = 1;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"BenchmarkFloatVlen","sequence":[{"name":"floatArray","type":{"stream":{"items":{"array":{"items":"float32","dimensions":2}}}}}]},"types":null}');
    end
  end

  methods (Abstract, Access=protected)
    write_float_array_(self, value)

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
        name = "write_float_array or end_float_array";
      else
        name = '<unknown>';
      end
    end
  end
end
