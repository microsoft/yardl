% This file was generated by the "yardl" tool. DO NOT EDIT.

% Abstract writer for protocol Scalars
classdef (Abstract) ScalarsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function self = ScalarsWriterBase()
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
    function write_int32(self, value)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.write_int32_(value);
      self.state_ = 1;
    end

    % Ordinal 1
    function write_record(self, value)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      self.write_record_(value);
      self.state_ = 2;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"Scalars","sequence":[{"name":"int32","type":"int32"},{"name":"record","type":"TestModel.RecordWithPrimitives"}]},"types":[{"name":"RecordWithPrimitives","fields":[{"name":"boolField","type":"bool"},{"name":"int8Field","type":"int8"},{"name":"uint8Field","type":"uint8"},{"name":"int16Field","type":"int16"},{"name":"uint16Field","type":"uint16"},{"name":"int32Field","type":"int32"},{"name":"uint32Field","type":"uint32"},{"name":"int64Field","type":"int64"},{"name":"uint64Field","type":"uint64"},{"name":"sizeField","type":"size"},{"name":"float32Field","type":"float32"},{"name":"float64Field","type":"float64"},{"name":"complexfloat32Field","type":"complexfloat32"},{"name":"complexfloat64Field","type":"complexfloat64"},{"name":"dateField","type":"date"},{"name":"timeField","type":"time"},{"name":"datetimeField","type":"datetime"}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_int32_(self, value)
    write_record_(self, value)

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
        name = "write_int32";
      elseif state == 1
        name = "write_record";
      else
        name = '<unknown>';
      end
    end
  end
end
