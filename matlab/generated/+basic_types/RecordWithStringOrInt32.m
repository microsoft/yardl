% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithStringOrInt32 < yardl.Union
  methods (Static)
    function res = RecordWithString(value)
      res = basic_types.RecordWithStringOrInt32(1, value);
    end

    function res = Int32(value)
      res = basic_types.RecordWithStringOrInt32(2, value);
    end

    function z = zeros(varargin)
      elem = basic_types.RecordWithStringOrInt32(0, yardl.None);
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
    function res = isRecordWithString(self)
      res = self.index == 1;
    end

    function res = isInt32(self)
      res = self.index == 2;
    end

    function eq = eq(self, other)
      eq = isa(other, "basic_types.RecordWithStringOrInt32") && all([self.index_] == [other.index_], 'all') && all([self.value] == [other.value], 'all');
    end

    function ne = ne(self, other)
      ne = ~self.eq(other);
    end

    function t = tag(self)
      tags_ = ["RecordWithString", "Int32"];
      t = tags_(self.index_);
    end
  end
end
