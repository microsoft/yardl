classdef ScalarOptionalsReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = ScalarOptionalsReaderBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 8
        if mod(obj.state_, 2) == 1
          previous_method = obj.state_to_method_name_(obj.state_ - 1);
          throw(yardl.ProtocolError("Protocol reader closed before all data was consumed. The iterable returned by '%s' was not fully consumed.", previous_method));
        else
          expected_method = obj.state_to_method_name_(obj.state_);
          throw(yardl.ProtocolError("Protocol reader closed before all data was consumed. Expected call to '%s'.", expected_method));
        end
      end
    end

    % Ordinal 0
    function value = read_optional_int(obj)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      value = obj.read_optional_int_();
      obj.state_ = 2;
    end

    % Ordinal 1
    function value = read_optional_record(obj)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      value = obj.read_optional_record_();
      obj.state_ = 4;
    end

    % Ordinal 2
    function value = read_record_with_optional_fields(obj)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      value = obj.read_record_with_optional_fields_();
      obj.state_ = 6;
    end

    % Ordinal 3
    function value = read_optional_record_with_optional_fields(obj)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      value = obj.read_optional_record_with_optional_fields_();
      obj.state_ = 8;
    end

    function copy_to(obj, writer)
      writer.write_optional_int(obj.read_optional_int());
      writer.write_optional_record(obj.read_optional_record());
      writer.write_record_with_optional_fields(obj.read_record_with_optional_fields());
      writer.write_optional_record_with_optional_fields(obj.read_optional_record_with_optional_fields());
    end
  end

  methods (Static)
    function res = schema()
      res = test_model.ScalarOptionalsWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    read_optional_int_(obj, value)
    read_optional_record_(obj, value)
    read_record_with_optional_fields_(obj, value)
    read_optional_record_with_optional_fields_(obj, value)

    close_(obj)
  end

  methods (Access=private)
    function raise_unexpected_state_(obj, actual)
      actual_method = obj.state_to_method_name_(actual);
      expected_method = obj.state_to_method_name_(obj.state_);
      throw(yardl.ProtocolError("Expected call to '%s' but received call to '%s'.", expected_method, actual_method));
    end

    function name = state_to_method_name_(obj, state)
      if state == 0
        name = 'read_optional_int';
      elseif state == 2
        name = 'read_optional_record';
      elseif state == 4
        name = 'read_record_with_optional_fields';
      elseif state == 6
        name = 'read_optional_record_with_optional_fields';
      else
        name = '<unknown>';
      end
    end
  end
end
