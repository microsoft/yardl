% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef GenericUnionWithRepeatedTypeParameters < yardl.Union
  methods (Static)
    function res = T(value)
      res = test_model.GenericUnionWithRepeatedTypeParameters(1, value);
    end

    function res = Tv(value)
      res = test_model.GenericUnionWithRepeatedTypeParameters(2, value);
    end

    function res = Ta(value)
      res = test_model.GenericUnionWithRepeatedTypeParameters(3, value);
    end

    function z = zeros(varargin)
      elem = test_model.GenericUnionWithRepeatedTypeParameters(0, yardl.None);
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
    function res = isT(self)
      res = self.index == 1;
    end

    function res = isTv(self)
      res = self.index == 2;
    end

    function res = isTa(self)
      res = self.index == 3;
    end

    function eq = eq(self, other)
      eq = isa(other, 'test_model.GenericUnionWithRepeatedTypeParameters') && other.index == self.index && other.value == self.value;
    end

    function ne = ne(self, other)
      ne = ~self.eq(other);
    end
  end
end
