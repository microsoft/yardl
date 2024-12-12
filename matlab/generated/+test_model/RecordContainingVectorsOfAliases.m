% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordContainingVectorsOfAliases < handle
  properties
    strings
    maps
    arrays
    tuples
  end

  methods
    function self = RecordContainingVectorsOfAliases(kwargs)
      arguments
        kwargs.strings = string.empty();
        kwargs.maps = test_model.AliasedMap.empty();
        kwargs.arrays = single.empty();
        kwargs.tuples = test_model.MyTuple.empty();
      end
      self.strings = kwargs.strings;
      self.maps = kwargs.maps;
      self.arrays = kwargs.arrays;
      self.tuples = kwargs.tuples;
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordContainingVectorsOfAliases") && ...
        isequal({self.strings}, {other.strings}) && ...
        isequal({self.maps}, {other.maps}) && ...
        isequal({self.arrays}, {other.arrays}) && ...
        isequal({self.tuples}, {other.tuples});
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
      elem = test_model.RecordContainingVectorsOfAliases();
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
