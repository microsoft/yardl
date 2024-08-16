% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef UnionOfContainerRecords < yardl.Union
  methods (Static)
    function res = RecordWithIntVectors(value)
      res = test_model.UnionOfContainerRecords(1, value);
    end

    function res = RecordWithFloatArrays(value)
      res = test_model.UnionOfContainerRecords(2, value);
    end

    function z = zeros(varargin)
      elem = test_model.UnionOfContainerRecords(0, yardl.None);
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
    function res = isRecordWithIntVectors(self)
      res = self.index == 1;
    end

    function res = isRecordWithFloatArrays(self)
      res = self.index == 2;
    end

    function eq = eq(self, other)
      eq = isa(other, "test_model.UnionOfContainerRecords") && isequal(self.index, other.index) && isequal(self.value, other.value);
    end

    function ne = ne(self, other)
      ne = ~self.eq(other);
    end

    function t = tag(self)
      tags_ = ["RecordWithIntVectors", "RecordWithFloatArrays"];
      t = tags_(self.index_);
    end
  end
end