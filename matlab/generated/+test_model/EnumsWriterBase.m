% Abstract writer for protocol Enums
classdef (Abstract) EnumsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = EnumsWriterBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 6
        expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_single(obj, value)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_single_(value);
      obj.state_ = 2;
    end

    % Ordinal 1
    function write_vec(obj, value)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      obj.write_vec_(value);
      obj.state_ = 4;
    end

    % Ordinal 2
    function write_size(obj, value)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      obj.write_size_(value);
      obj.state_ = 6;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"Enums","sequence":[{"name":"single","type":"TestModel.Fruits"},{"name":"vec","type":{"vector":{"items":"TestModel.Fruits"}}},{"name":"size","type":"TestModel.SizeBasedEnum"}]},"types":[{"name":"Fruits","values":[{"symbol":"apple","value":0},{"symbol":"banana","value":1},{"symbol":"pear","value":2}]},{"name":"Fruits","type":"BasicTypes.Fruits"},{"name":"SizeBasedEnum","base":"size","values":[{"symbol":"a","value":0},{"symbol":"b","value":1},{"symbol":"c","value":2}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_single_(obj, value)
    write_vec_(obj, value)
    write_size_(obj, value)

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
        name = 'write_single';
      elseif state == 2
        name = 'write_vec';
      elseif state == 4
        name = 'write_size';
      else
        name = '<unknown>';
      end
    end
  end
end