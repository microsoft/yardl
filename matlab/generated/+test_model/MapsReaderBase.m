% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MapsReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function self = MapsReaderBase()
      self.state_ = 0;
    end

    function close(self)
      self.close_();
      if self.state_ ~= 4
        expected_method = self.state_to_method_name_(self.state_);
        throw(yardl.ProtocolError("Protocol reader closed before all data was consumed. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function value = read_string_to_int(self)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      value = self.read_string_to_int_();
      self.state_ = 1;
    end

    % Ordinal 1
    function value = read_int_to_string(self)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      value = self.read_int_to_string_();
      self.state_ = 2;
    end

    % Ordinal 2
    function value = read_string_to_union(self)
      if self.state_ ~= 2
        self.raise_unexpected_state_(2);
      end

      value = self.read_string_to_union_();
      self.state_ = 3;
    end

    % Ordinal 3
    function value = read_aliased_generic(self)
      if self.state_ ~= 3
        self.raise_unexpected_state_(3);
      end

      value = self.read_aliased_generic_();
      self.state_ = 4;
    end

    function copy_to(self, writer)
      writer.write_string_to_int(self.read_string_to_int());
      writer.write_int_to_string(self.read_int_to_string());
      writer.write_string_to_union(self.read_string_to_union());
      writer.write_aliased_generic(self.read_aliased_generic());
    end
  end

  methods (Static)
    function res = schema()
      res = test_model.MapsWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    read_string_to_int_(self)
    read_int_to_string_(self)
    read_string_to_union_(self)
    read_aliased_generic_(self)

    close_(self)
  end

  methods (Access=private)
    function raise_unexpected_state_(self, actual)
      actual_method = self.state_to_method_name_(actual);
      expected_method = self.state_to_method_name_(self.state_);
      throw(yardl.ProtocolError("Expected call to '%s' but received call to '%s'.", expected_method, actual_method));
    end

    function name = state_to_method_name_(self, state)
      if state == 0
        name = "read_string_to_int";
      elseif state == 1
        name = "read_int_to_string";
      elseif state == 2
        name = "read_string_to_union";
      elseif state == 3
        name = "read_aliased_generic";
      else
        name = "<unknown>";
      end
    end
  end
end
