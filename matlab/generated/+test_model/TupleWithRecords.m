% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef TupleWithRecords < handle
  properties
    a
    b
  end

  methods
    function self = TupleWithRecords(a, b)
      if nargin > 0
        self.a = a;
        self.b = b;
      else
        self.a = test_model.SimpleRecord();
        self.b = test_model.SimpleRecord();
      end
    end

    function res = eq(self, other)
      res = ...
        isa(other, 'test_model.TupleWithRecords') && ...
        all([self.a] == [other.a]) && ...
        all([self.b] == [other.b]);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.TupleWithRecords();
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
