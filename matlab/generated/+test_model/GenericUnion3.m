% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef GenericUnion3 < yardl.Union
  methods (Static)
    function res = T(value)
      res = test_model.GenericUnion3(1, value);
    end

    function res = U(value)
      res = test_model.GenericUnion3(2, value);
    end

    function res = V(value)
      res = test_model.GenericUnion3(3, value);
    end

    function z = zeros(varargin)
      elem = test_model.GenericUnion3(0, yardl.None);
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
    function eq = eq(self, other)
      eq = isa(other, 'test_model.GenericUnion3') && other.index == self.index && other.value == self.value;
    end

    function ne = ne(self, other)
      ne = ~self.eq(other);
    end
  end
end
