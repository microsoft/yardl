% This file was generated by the "yardl" tool. DO NOT EDIT.

% Abstract writer for protocol StreamsOfUnions
classdef (Abstract) StreamsOfUnionsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function self = StreamsOfUnionsWriterBase()
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
    function write_int_or_simple_record(self, value)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.write_int_or_simple_record_(value);
    end

    function end_int_or_simple_record(self)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.end_stream_();
      self.state_ = 1;
    end

    % Ordinal 1
    function write_nullable_int_or_simple_record(self, value)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      self.write_nullable_int_or_simple_record_(value);
    end

    function end_nullable_int_or_simple_record(self)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      self.end_stream_();
      self.state_ = 2;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"StreamsOfUnions","sequence":[{"name":"intOrSimpleRecord","type":{"stream":{"items":[{"tag":"int32","type":"int32"},{"tag":"SimpleRecord","type":"TestModel.SimpleRecord"}]}}},{"name":"nullableIntOrSimpleRecord","type":{"stream":{"items":[null,{"tag":"int32","type":"int32"},{"tag":"SimpleRecord","type":"TestModel.SimpleRecord"}]}}}]},"types":[{"name":"SimpleRecord","fields":[{"name":"x","type":"int32"},{"name":"y","type":"int32"},{"name":"z","type":"int32"}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_int_or_simple_record_(self, value)
    write_nullable_int_or_simple_record_(self, value)

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
        name = "write_int_or_simple_record or end_int_or_simple_record";
      elseif state == 1
        name = "write_nullable_int_or_simple_record or end_nullable_int_or_simple_record";
      else
        name = '<unknown>';
      end
    end
  end
end
