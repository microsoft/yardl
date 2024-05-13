% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef AliasedNullableIntSimpleRecord < yardl.Union
  methods (Static)
    function res = Int32(value)
      res = test_model.AliasedNullableIntSimpleRecord(1, value);
    end

    function res = SimpleRecord(value)
      res = test_model.AliasedNullableIntSimpleRecord(2, value);
    end

    function z = zeros(varargin)
      elem = test_model.AliasedNullableIntSimpleRecord(0, yardl.None);
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
    function res = isInt32(self)
      res = self.index == 1;
    end

    function res = isSimpleRecord(self)
      res = self.index == 2;
    end

    function eq = eq(self, other)
      eq = isa(other, "test_model.AliasedNullableIntSimpleRecord") && other.index == self.index && all([self.value] == [other.value]);
    end

    function ne = ne(self, other)
      ne = ~self.eq(other);
    end

    function t = tag(self)
      tags_ = ["Int32", "SimpleRecord"];
      t = tags_(self.index_);
    end
  end
end
