% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithVlens < handle
  properties
    a
    b
    c
  end

  methods
    function self = RecordWithVlens(kwargs)
      arguments
        kwargs.a = test_model.SimpleRecord.empty();
        kwargs.b = int32(0);
        kwargs.c = int32(0);
      end
      self.a = kwargs.a;
      self.b = kwargs.b;
      self.c = kwargs.c;
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordWithVlens") && ...
        isequal({self.a}, {other.a}) && ...
        isequal({self.b}, {other.b}) && ...
        isequal({self.c}, {other.c});
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
