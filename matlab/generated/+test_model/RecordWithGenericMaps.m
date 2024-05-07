% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithGenericMaps < handle
  properties
    m
    am
  end

  methods
    function self = RecordWithGenericMaps(kwargs)
      arguments
        kwargs.m = dictionary;
        kwargs.am = dictionary;
      end
      self.m = kwargs.m;
      self.am = kwargs.am;
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordWithGenericMaps") && ...
        isequal(self.m, other.m) && ...
        isequal(self.am, other.am);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithGenericMaps();
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
