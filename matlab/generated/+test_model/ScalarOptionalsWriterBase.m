% Abstract writer for protocol ScalarOptionals
classdef (Abstract) ScalarOptionalsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = ScalarOptionalsWriterBase()
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
    function write_optional_int(obj, value)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_optional_int_(value);
      obj.state_ = 2;
    end

    % Ordinal 1
    function write_optional_record(obj, value)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      obj.write_optional_record_(value);
      obj.state_ = 4;
    end

    % Ordinal 2
    function write_record_with_optional_fields(obj, value)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      obj.write_record_with_optional_fields_(value);
      obj.state_ = 6;
    end

    % Ordinal 3
    function write_optional_record_with_optional_fields(obj, value)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      obj.write_optional_record_with_optional_fields_(value);
      obj.state_ = 8;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"ScalarOptionals","sequence":[{"name":"optionalInt","type":[null,"int32"]},{"name":"optionalRecord","type":[null,"TestModel.SimpleRecord"]},{"name":"recordWithOptionalFields","type":"TestModel.RecordWithOptionalFields"},{"name":"optionalRecordWithOptionalFields","type":[null,"TestModel.RecordWithOptionalFields"]}]},"types":[{"name":"RecordWithOptionalFields","fields":[{"name":"optionalInt","type":[null,"int32"]},{"name":"optionalIntAlternateSyntax","type":[null,"int32"]},{"name":"optionalTime","type":[null,"time"]}]},{"name":"SimpleRecord","fields":[{"name":"x","type":"int32"},{"name":"y","type":"int32"},{"name":"z","type":"int32"}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_optional_int_(obj, value)
    write_optional_record_(obj, value)
    write_record_with_optional_fields_(obj, value)
    write_optional_record_with_optional_fields_(obj, value)

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
        name = 'write_optional_int';
      elseif state == 2
        name = 'write_optional_record';
      elseif state == 4
        name = 'write_record_with_optional_fields';
      elseif state == 6
        name = 'write_optional_record_with_optional_fields';
      else
        name = '<unknown>';
      end
    end
  end
end
