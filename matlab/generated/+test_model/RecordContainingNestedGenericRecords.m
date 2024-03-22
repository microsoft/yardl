classdef RecordContainingNestedGenericRecords < handle
  properties
    f1
    f1a
    f2
    f2a
    nested
  end

  methods
    function obj = RecordContainingNestedGenericRecords(f1, f1a, f2, f2a, nested)
      if nargin > 0
        obj.f1 = f1;
        obj.f1a = f1a;
        obj.f2 = f2;
        obj.f2a = f2a;
        obj.nested = nested;
      else
        obj.f1 = test_model.RecordWithOptionalGenericField(yardl.None);
        obj.f1a = test_model.RecordWithAliasedOptionalGenericField(yardl.None);
        obj.f2 = test_model.RecordWithOptionalGenericUnionField(yardl.None);
        obj.f2a = test_model.RecordWithAliasedOptionalGenericUnionField(yardl.None);
        obj.nested = test_model.RecordContainingGenericRecords(test_model.RecordWithOptionalGenericField(yardl.None), test_model.RecordWithAliasedOptionalGenericField(yardl.None), test_model.RecordWithOptionalGenericUnionField(yardl.None), test_model.RecordWithAliasedOptionalGenericUnionField(yardl.None), tuples.Tuple("", int32(0)), tuples.Tuple("", int32(0)), test_model.RecordWithGenericVectors(int32.empty(), int32.empty()), test_model.RecordWithGenericFixedVectors(repelem(int32(0), 3), repelem(int32(0), 3)), test_model.RecordWithGenericArrays(int32.empty(0, 0), repelem(int32(0), 8, 16), int32.empty(), int32.empty(0, 0), repelem(int32(0), 8, 16), int32.empty()), test_model.RecordWithGenericMaps(dictionary, dictionary));
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordContainingNestedGenericRecords') && ...
        all([obj.f1] == [other.f1]) && ...
        isequal(obj.f1a, other.f1a) && ...
        all([obj.f2] == [other.f2]) && ...
        isequal(obj.f2a, other.f2a) && ...
        isequal(obj.nested, other.nested);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordContainingNestedGenericRecords();
      if nargin == 0
        z = elem;
      elseif nargin == 1
        n = varargin{1};
        z = reshape(repelem(elem, n*n), [n, n]);
      else
        sz = [varargin{:}];
        z = reshape(repelem(elem, prod(sz)), sz);
      end
    end
  end
end