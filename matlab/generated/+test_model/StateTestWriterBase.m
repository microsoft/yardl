% This file was generated by the "yardl" tool. DO NOT EDIT.

% Abstract writer for protocol StateTest
classdef (Abstract) StateTestWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function self = StateTestWriterBase()
      self.state_ = 0;
    end

    function close(self)
      self.close_();
      if self.state_ ~= 3
        expected_method = self.state_to_method_name_(self.state_);
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_an_int(self, value)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.write_an_int_(value);
      self.state_ = 1;
    end

    % Ordinal 1
    function write_a_stream(self, value)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      self.write_a_stream_(value);
    end

    function end_a_stream(self)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      self.end_stream_();
      self.state_ = 2;
    end

    % Ordinal 2
    function write_another_int(self, value)
      if self.state_ ~= 2
        self.raise_unexpected_state_(2);
      end

      self.write_another_int_(value);
      self.state_ = 3;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"StateTest","sequence":[{"name":"anInt","type":"int32"},{"name":"aStream","type":{"stream":{"items":"int32"}}},{"name":"anotherInt","type":"int32"}]},"types":null}');
    end
  end

  methods (Abstract, Access=protected)
    write_an_int_(self, value)
    write_a_stream_(self, value)
    write_another_int_(self, value)

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
        name = "write_an_int";
      elseif state == 1
        name = "write_a_stream or end_a_stream";
      elseif state == 2
        name = "write_another_int";
      else
        name = '<unknown>';
      end
    end
  end
end
