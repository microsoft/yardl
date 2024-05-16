% This file was generated by the "yardl" tool. DO NOT EDIT.

% Abstract writer for protocol Enums
classdef (Abstract) EnumsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function self = EnumsWriterBase()
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
    function write_single(self, value)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.write_single_(value);
      self.state_ = 1;
    end

    % Ordinal 1
    function write_vec(self, value)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      self.write_vec_(value);
      self.state_ = 2;
    end

    % Ordinal 2
    function write_size(self, value)
      if self.state_ ~= 2
        self.raise_unexpected_state_(2);
      end

      self.write_size_(value);
      self.state_ = 3;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"Enums","sequence":[{"name":"single","type":"TestModel.Fruits"},{"name":"vec","type":{"vector":{"items":"TestModel.Fruits"}}},{"name":"size","type":"TestModel.SizeBasedEnum"}]},"types":[{"name":"Fruits","values":[{"symbol":"apple","value":0},{"symbol":"banana","value":1},{"symbol":"pear","value":2}]},{"name":"Fruits","type":"BasicTypes.Fruits"},{"name":"SizeBasedEnum","base":"size","values":[{"symbol":"a","value":0},{"symbol":"b","value":1},{"symbol":"c","value":2}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_single_(self, value)
    write_vec_(self, value)
    write_size_(self, value)

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
        name = "write_single";
      elseif state == 1
        name = "write_vec";
      elseif state == 2
        name = "write_size";
      else
        name = '<unknown>';
      end
    end
  end
end
