% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef AliasedMultiGenericOptional < yardl.Union
  methods (Static)
    function res = T(value)
      res = test_model.AliasedMultiGenericOptional(1, value);
    end

    function res = U(value)
      res = test_model.AliasedMultiGenericOptional(2, value);
    end

    function z = zeros(varargin)
      elem = test_model.AliasedMultiGenericOptional(0, yardl.None);
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

    function res = isU(self)
      res = self.index == 2;
    end

    function eq = eq(self, other)
      eq = isa(other, "test_model.AliasedMultiGenericOptional") && all([self.index_] == [other.index_], 'all') && all([self.value] == [other.value], 'all');
    end

    function ne = ne(self, other)
      ne = ~self.eq(other);
    end

    function t = tag(self)
      tags_ = ["T", "U"];
      t = tags_(self.index_);
    end
  end
end
