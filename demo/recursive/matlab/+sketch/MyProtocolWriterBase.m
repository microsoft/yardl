% This file was generated by the "yardl" tool. DO NOT EDIT.

% Abstract writer for protocol MyProtocol
classdef (Abstract) MyProtocolWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function self = MyProtocolWriterBase()
      self.state_ = 0;
    end

    function close(self)
      self.close_();
      if self.state_ ~= 4
        expected_method = self.state_to_method_name_(self.state_);
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_tree(self, value)
      if self.state_ ~= 0
        self.raise_unexpected_state_(0);
      end

      self.write_tree_(value);
      self.state_ = 1;
    end

    % Ordinal 1
    function write_ptree(self, value)
      if self.state_ ~= 1
        self.raise_unexpected_state_(1);
      end

      self.write_ptree_(value);
      self.state_ = 2;
    end

    % Ordinal 2
    function write_list(self, value)
      if self.state_ ~= 2
        self.raise_unexpected_state_(2);
      end

      self.write_list_(value);
      self.state_ = 3;
    end

    % Ordinal 3
    % dirs: !stream
    %   items: Directory
    function write_cwd(self, value)
      if self.state_ ~= 3
        self.raise_unexpected_state_(3);
      end

      self.write_cwd_(value);
    end

    function end_cwd(self)
      if self.state_ ~= 3
        self.raise_unexpected_state_(3);
      end

      self.end_stream_();
      self.state_ = 4;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"MyProtocol","sequence":[{"name":"tree","type":"Sketch.BinaryTree"},{"name":"ptree","type":"Sketch.BinaryTree"},{"name":"list","type":[null,{"name":"Sketch.LinkedList","typeArguments":["string"]}]},{"name":"cwd","type":{"stream":{"items":"Sketch.DirectoryEntry"}}}]},"types":[{"name":"BinaryTree","fields":[{"name":"value","type":"int32"},{"name":"left","type":"Sketch.BinaryTree"},{"name":"right","type":"Sketch.BinaryTree"}]},{"name":"Directory","fields":[{"name":"name","type":"string"},{"name":"entries","type":{"vector":{"items":"Sketch.DirectoryEntry"}}}]},{"name":"DirectoryEntry","type":[{"tag":"File","type":"Sketch.File"},{"tag":"Directory","type":"Sketch.Directory"}]},{"name":"File","fields":[{"name":"name","type":"string"},{"name":"data","type":{"vector":{"items":"uint8"}}}]},{"name":"LinkedList","typeParameters":["T"],"fields":[{"name":"value","type":"T"},{"name":"next","type":{"name":"Sketch.LinkedList","typeArguments":["T"]}}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_tree_(self, value)
    write_ptree_(self, value)
    write_list_(self, value)
    write_cwd_(self, value)

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
        name = "write_tree";
      elseif state == 1
        name = "write_ptree";
      elseif state == 2
        name = "write_list";
      elseif state == 3
        name = "write_cwd or end_cwd";
      else
        name = '<unknown>';
      end
    end
  end
end