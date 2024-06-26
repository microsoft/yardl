% This file was generated by the "yardl" tool. DO NOT EDIT.

% Abstract writer for protocol Streams
classdef (Abstract) StreamsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function self = StreamsWriterBase()
      self.state_ = 0;
    end

    function close(self)
      self.close_();
      if self.state_ ~= 4
        expected_method = self.state_to_method_name_(self.state_);
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_int_data(self, value)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.write_int_data_(value);
    end

    function end_int_data(self)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.end_stream_();
      self.state_ = 1;
    end

    % Ordinal 1
    function write_optional_int_data(self, value)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      self.write_optional_int_data_(value);
    end

    function end_optional_int_data(self)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      self.end_stream_();
      self.state_ = 2;
    end

    % Ordinal 2
    function write_record_with_optional_vector_data(self, value)
      if self.state_ ~= 2
        self.raise_unexpected_state_(2);
      end

      self.write_record_with_optional_vector_data_(value);
    end

    function end_record_with_optional_vector_data(self)
      if self.state_ ~= 2
        self.raise_unexpected_state_(2);
      end

      self.end_stream_();
      self.state_ = 3;
    end

    % Ordinal 3
    function write_fixed_vector(self, value)
      if self.state_ ~= 3
        self.raise_unexpected_state_(3);
      end

      self.write_fixed_vector_(value);
    end

    function end_fixed_vector(self)
      if self.state_ ~= 3
        self.raise_unexpected_state_(3);
      end

      self.end_stream_();
      self.state_ = 4;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"Streams","sequence":[{"name":"intData","type":{"stream":{"items":"int32"}}},{"name":"optionalIntData","type":{"stream":{"items":[null,"int32"]}}},{"name":"recordWithOptionalVectorData","type":{"stream":{"items":"TestModel.RecordWithOptionalVector"}}},{"name":"fixedVector","type":{"stream":{"items":{"vector":{"items":"int32","length":3}}}}}]},"types":[{"name":"RecordWithOptionalVector","fields":[{"name":"optionalVector","type":[null,{"vector":{"items":"int32"}}]}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_int_data_(self, value)
    write_optional_int_data_(self, value)
    write_record_with_optional_vector_data_(self, value)
    write_fixed_vector_(self, value)

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
        name = "write_int_data or end_int_data";
      elseif state == 1
        name = "write_optional_int_data or end_optional_int_data";
      elseif state == 2
        name = "write_record_with_optional_vector_data or end_record_with_optional_vector_data";
      elseif state == 3
        name = "write_fixed_vector or end_fixed_vector";
      else
        name = '<unknown>';
      end
    end
  end
end
