% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef File < handle
  properties
    name
    data
  end

  methods
    function self = File(kwargs)
      arguments
        kwargs.name = "";
        kwargs.data = uint8.empty();
      end
      self.name = kwargs.name;
      self.data = kwargs.data;
    end

    function res = eq(self, other)
      res = ...
        isa(other, "sketch.File") && ...
        isequal(self.name, other.name) && ...
        isequal(self.data, other.data);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = sketch.File();
      if nargin == 0
        z = elem;
        return;
      end
      sz = [varargin{:}];
      if isscalar(sz)
        sz = [sz, sz];
      end
      z = reshape(repelem(elem, prod(sz)), sz);
    end
  end
end