% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef SimpleRecord < handle
  properties
    x
    y
    z
  end

  methods
    function self = SimpleRecord(x, y, z)
      if nargin > 0
        self.x = x;
        self.y = y;
        self.z = z;
      else
        self.x = int32(0);
        self.y = int32(0);
        self.z = int32(0);
      end
    end

    function res = eq(self, other)
      res = ...
        isa(other, 'test_model.SimpleRecord') && ...
        all([self.x] == [other.x]) && ...
        all([self.y] == [other.y]) && ...
        all([self.z] == [other.z]);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.SimpleRecord();
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
