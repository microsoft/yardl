% This file was generated by the "yardl" tool. DO NOT EDIT.

% Abstract writer for protocol MultiDArrays
classdef (Abstract) MultiDArraysWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function self = MultiDArraysWriterBase()
      self.state_ = 0;
    end

    function close(self)
      self.close_();
      if self.state_ ~= 2
        expected_method = self.state_to_method_name_(self.state_);
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_images(self, value)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.write_images_(value);
    end

    function end_images(self)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.end_stream_();
      self.state_ = 1;
    end

    % Ordinal 1
    function write_frames(self, value)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      self.write_frames_(value);
    end

    function end_frames(self)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      self.end_stream_();
      self.state_ = 2;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"MultiDArrays","sequence":[{"name":"images","type":{"stream":{"items":{"array":{"items":"float32","dimensions":[{"name":"ch"},{"name":"z"},{"name":"y"},{"name":"x"}]}}}}},{"name":"frames","type":{"stream":{"items":{"array":{"items":"float32","dimensions":[{"name":"ch","length":1},{"name":"z","length":1},{"name":"y","length":64},{"name":"x","length":32}]}}}}}]},"types":null}');
    end
  end

  methods (Abstract, Access=protected)
    write_images_(self, value)
    write_frames_(self, value)

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
        name = "write_images or end_images";
      elseif state == 1
        name = "write_frames or end_frames";
      else
        name = '<unknown>';
      end
    end
  end
end