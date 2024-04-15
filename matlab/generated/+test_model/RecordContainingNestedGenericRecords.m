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
    function self = RecordContainingNestedGenericRecords(f1, f1a, f2, f2a, nested)
      if nargin > 0
        self.f1 = f1;
        self.f1a = f1a;
        self.f2 = f2;
        self.f2a = f2a;
        self.nested = nested;
      else
        self.f1 = test_model.RecordWithOptionalGenericField(yardl.None);
        self.f1a = test_model.RecordWithAliasedOptionalGenericField(yardl.None);
        self.f2 = test_model.RecordWithOptionalGenericUnionField(yardl.None);
        self.f2a = test_model.RecordWithAliasedOptionalGenericUnionField(yardl.None);
        self.nested = test_model.RecordContainingGenericRecords(test_model.RecordWithOptionalGenericField(yardl.None), test_model.RecordWithAliasedOptionalGenericField(yardl.None), test_model.RecordWithOptionalGenericUnionField(yardl.None), test_model.RecordWithAliasedOptionalGenericUnionField(yardl.None), tuples.Tuple("", int32(0)), tuples.Tuple("", int32(0)), test_model.RecordWithGenericVectors(int32.empty(), int32.empty()), test_model.RecordWithGenericFixedVectors(repelem(int32(0), 3), repelem(int32(0), 3)), test_model.RecordWithGenericArrays(int32.empty(0, 0), repelem(int32(0), 8, 16), int32.empty(), int32.empty(0, 0), repelem(int32(0), 8, 16), int32.empty()), test_model.RecordWithGenericMaps(dictionary, dictionary));
      end
    end

    function res = eq(self, other)
      res = ...
        isa(other, 'test_model.RecordContainingNestedGenericRecords') && ...
        all([self.f1] == [other.f1]) && ...
        isequal(self.f1a, other.f1a) && ...
        all([self.f2] == [other.f2]) && ...
        isequal(self.f2a, other.f2a) && ...
        isequal(self.nested, other.nested);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
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
