% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MapOrScalar < yardl.Union
  methods (Static)
    function res = Map(value)
      res = test_model.MapOrScalar(1, value);
    end

    function res = Scalar(value)
      res = test_model.MapOrScalar(2, value);
    end

    function z = zeros(varargin)
      elem = test_model.MapOrScalar(0, yardl.None);
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

  methods
    function res = isMap(self)
      res = self.index == 1;
    end

    function res = isScalar(self)
      res = self.index == 2;
    end

    function eq = eq(self, other)
      eq = isa(other, "test_model.MapOrScalar") && other.index == self.index && other.value == self.value;
    end

    function ne = ne(self, other)
      ne = ~self.eq(other);
    end
  end
end
