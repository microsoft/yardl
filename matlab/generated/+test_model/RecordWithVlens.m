% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithVlens < handle
  properties
    a
    b
    c
  end

  methods
    function self = RecordWithVlens(a, b, c)
      if nargin > 0
        self.a = a;
        self.b = b;
        self.c = c;
      else
        self.a = test_model.SimpleRecord.empty();
        self.b = int32(0);
        self.c = int32(0);
      end
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordWithVlens") && ...
        all([self.a] == [other.a]) && ...
        all([self.b] == [other.b]) && ...
        all([self.c] == [other.c]);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithVlens();
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
