% Abstract writer for protocol Streams
classdef (Abstract) StreamsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = StreamsWriterBase()
      obj.state_ = 0;
    end

    function close(obj)
      if obj.state_ == 7
        obj.end_stream_();
        obj.close_();
        return
      end
      obj.close_();
      if obj.state_ ~= 8
        expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_int_data(obj, value)
      if bitand(int32(obj.state_), bitcmp(1, 'int8')) ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_int_data_(value);
      obj.state_ = 1;
    end

    % Ordinal 1
    function write_optional_int_data(obj, value)
      if obj.state_ == 1
        obj.end_stream_();
        obj.state_ = 2;
      elseif bitand(int32(obj.state_), bitcmp(1, 'int8')) ~= 2
        obj.raise_unexpected_state_(2);
      end

      obj.write_optional_int_data_(value);
      obj.state_ = 3;
    end

    % Ordinal 2
    function write_record_with_optional_vector_data(obj, value)
      if obj.state_ == 3
        obj.end_stream_();
        obj.state_ = 4;
      elseif bitand(int32(obj.state_), bitcmp(1, 'int8')) ~= 4
        obj.raise_unexpected_state_(4);
      end

      obj.write_record_with_optional_vector_data_(value);
      obj.state_ = 5;
    end

    % Ordinal 3
    function write_fixed_vector(obj, value)
      if obj.state_ == 5
        obj.end_stream_();
        obj.state_ = 6;
      elseif bitand(int32(obj.state_), bitcmp(1, 'int8')) ~= 6
        obj.raise_unexpected_state_(6);
      end

      obj.write_fixed_vector_(value);
      obj.state_ = 7;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"Streams","sequence":[{"name":"intData","type":{"stream":{"items":"int32"}}},{"name":"optionalIntData","type":{"stream":{"items":[null,"int32"]}}},{"name":"recordWithOptionalVectorData","type":{"stream":{"items":"TestModel.RecordWithOptionalVector"}}},{"name":"fixedVector","type":{"stream":{"items":{"vector":{"items":"int32","length":3}}}}}]},"types":[{"name":"RecordWithOptionalVector","fields":[{"name":"optionalVector","type":[null,{"vector":{"items":"int32"}}]}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_int_data_(obj, value)
    write_optional_int_data_(obj, value)
    write_record_with_optional_vector_data_(obj, value)
    write_fixed_vector_(obj, value)

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
        name = 'write_int_data';
      elseif state == 2
        name = 'write_optional_int_data';
      elseif state == 4
        name = 'write_record_with_optional_vector_data';
      elseif state == 6
        name = 'write_fixed_vector';
      else
        name = '<unknown>';
      end
    end
  end
end
