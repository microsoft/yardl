% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithAliasedOptionalGenericUnionField < handle
  properties
    v
  end

  methods
    function self = RecordWithAliasedOptionalGenericUnionField(v)
      if nargin > 0
        self.v = v;
      else
        self.v = yardl.None;
      end
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordWithAliasedOptionalGenericUnionField") && ...
        isequal(self.v, other.v);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
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
