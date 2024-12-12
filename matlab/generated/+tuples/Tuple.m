% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef Tuple < handle
  properties
    v1
    v2
  end

  methods
    function self = Tuple(kwargs)
      arguments
        kwargs.v1;
        kwargs.v2;
      end
      if ~isfield(kwargs, "v1")
        throw(yardl.TypeError("Missing required keyword argument 'v1'"))
      end
      self.v1 = kwargs.v1;
      if ~isfield(kwargs, "v2")
        throw(yardl.TypeError("Missing required keyword argument 'v2'"))
      end
      self.v2 = kwargs.v2;
    end

    function res = eq(self, other)
      res = ...
        isa(other, "tuples.Tuple") && ...
        isequal({self.v1}, {other.v1}) && ...
        isequal({self.v2}, {other.v2});
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end

    function res = isequal(self, other)
      res = all(eq(self, other));
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = tuples.Tuple(v1=yardl.None, v2=yardl.None);
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
