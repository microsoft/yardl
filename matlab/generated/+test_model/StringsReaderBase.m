% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef StringsReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function self = StringsReaderBase()
      self.state_ = 0;
    end

    function close(self)
      self.close_();
      if self.state_ ~= 2
        expected_method = self.state_to_method_name_(self.state_);
        throw(yardl.ProtocolError("Protocol reader closed before all data was consumed. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function value = read_single_string(self)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      value = self.read_single_string_();
      self.state_ = 1;
    end

    % Ordinal 1
    function value = read_rec_with_string(self)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      value = self.read_rec_with_string_();
      self.state_ = 2;
    end

    function copy_to(self, writer)
      writer.write_single_string(self.read_single_string());
      writer.write_rec_with_string(self.read_rec_with_string());
    end
  end

  methods (Static)
    function res = schema()
      res = test_model.StringsWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    read_single_string_(self)
    read_rec_with_string_(self)

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
        name = "read_single_string";
      elseif state == 1
        name = "read_rec_with_string";
      else
        name = "<unknown>";
      end
    end
  end
end
