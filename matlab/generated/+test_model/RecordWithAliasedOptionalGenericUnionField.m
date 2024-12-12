% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithAliasedOptionalGenericUnionField < handle
  properties
    v
  end

  methods
    function self = RecordWithAliasedOptionalGenericUnionField(kwargs)
      arguments
        kwargs.v = yardl.None;
      end
      self.v = kwargs.v;
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordWithAliasedOptionalGenericUnionField") && ...
        isequal({self.v}, {other.v});
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
      elem = test_model.RecordWithAliasedOptionalGenericUnionField();
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
