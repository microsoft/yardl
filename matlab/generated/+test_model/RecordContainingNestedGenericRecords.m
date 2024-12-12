% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordContainingNestedGenericRecords < handle
  properties
    f1
    f1a
    f2
    f2a
    nested
  end

  methods
    function self = RecordContainingNestedGenericRecords(kwargs)
      arguments
        kwargs.f1 = test_model.RecordWithOptionalGenericField();
        kwargs.f1a = test_model.RecordWithAliasedOptionalGenericField();
        kwargs.f2 = test_model.RecordWithOptionalGenericUnionField();
        kwargs.f2a = test_model.RecordWithAliasedOptionalGenericUnionField();
        kwargs.nested = test_model.RecordContainingGenericRecords(g3=tuples.Tuple(v1="", v2=int32(0)), g3a=tuples.Tuple(v1="", v2=int32(0)), g4=test_model.RecordWithGenericVectors(v=int32.empty(), av=int32.empty()), g5=test_model.RecordWithGenericFixedVectors(fv=repelem(int32(0), 3), afv=repelem(int32(0), 3)), g6=test_model.RecordWithGenericArrays(nd=int32.empty(), fixed_nd=repelem(int32(0), 8, 16), dynamic_nd=int32.empty(), aliased_nd=int32.empty(), aliased_fixed_nd=repelem(int32(0), 8, 16), aliased_dynamic_nd=int32.empty()));
      end
      self.f1 = kwargs.f1;
      self.f1a = kwargs.f1a;
      self.f2 = kwargs.f2;
      self.f2a = kwargs.f2a;
      self.nested = kwargs.nested;
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordContainingNestedGenericRecords") && ...
        isequal({self.f1}, {other.f1}) && ...
        isequal({self.f1a}, {other.f1a}) && ...
        isequal({self.f2}, {other.f2}) && ...
        isequal({self.f2a}, {other.f2a}) && ...
        isequal({self.nested}, {other.nested});
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
      elem = test_model.RecordContainingNestedGenericRecords();
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
