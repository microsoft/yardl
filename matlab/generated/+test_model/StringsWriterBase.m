% Abstract writer for protocol Strings
classdef (Abstract) StringsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = StringsWriterBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 4
        expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_single_string(obj, value)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_single_string_(value);
      obj.state_ = 2;
    end

    % Ordinal 1
    function write_rec_with_string(obj, value)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      obj.write_rec_with_string_(value);
      obj.state_ = 4;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"Strings","sequence":[{"name":"singleString","type":"string"},{"name":"recWithString","type":"TestModel.RecordWithStrings"}]},"types":[{"name":"RecordWithStrings","fields":[{"name":"a","type":"string"},{"name":"b","type":"string"}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_single_string_(obj, value)
    write_rec_with_string_(obj, value)

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
        name = 'write_single_string';
      elseif state == 2
        name = 'write_rec_with_string';
      else
        name = '<unknown>';
      end
    end
  end
end
