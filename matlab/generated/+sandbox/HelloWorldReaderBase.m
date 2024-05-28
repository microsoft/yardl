% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef HelloWorldReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function self = HelloWorldReaderBase()
      self.state_ = 0;
    end

    function close(self)
      self.close_();
      if self.state_ ~= 1
        expected_method = self.state_to_method_name_(self.state_);
        throw(yardl.ProtocolError("Protocol reader closed before all data was consumed. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function more = has_data(self)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      more = self.has_data_();
      if ~more
        self.state_ = 1;
      end
    end

    function value = read_data(self)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      value = self.read_data_();
    end

    function copy_to(self, writer)
      while self.has_data()
        item = self.read_data();
        writer.write_data({item});
      end
      writer.end_data();
    end
  end

  methods (Static)
    function res = schema()
      res = sandbox.HelloWorldWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    has_data_(self)
    read_data_(self)

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
        name = "read_data";
      else
        name = "<unknown>";
      end
    end
  end
end