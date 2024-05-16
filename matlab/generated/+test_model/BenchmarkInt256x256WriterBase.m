% This file was generated by the "yardl" tool. DO NOT EDIT.

% Abstract writer for protocol BenchmarkInt256x256
classdef (Abstract) BenchmarkInt256x256WriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function self = BenchmarkInt256x256WriterBase()
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
    function write_int256x256(self, value)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.write_int256x256_(value);
    end

    function end_int256x256(self)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.end_stream_();
      self.state_ = 1;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"BenchmarkInt256x256","sequence":[{"name":"int256x256","type":{"stream":{"items":{"array":{"items":"int32","dimensions":[{"length":256},{"length":256}]}}}}}]},"types":null}');
    end
  end

  methods (Abstract, Access=protected)
    write_int256x256_(self, value)

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
        name = "write_int256x256 or end_int256x256";
      else
        name = '<unknown>';
      end
    end
  end
end
