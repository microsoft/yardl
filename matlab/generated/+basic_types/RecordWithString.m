% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithString < handle
  properties
    i
  end

  methods
    function self = RecordWithString(kwargs)
      arguments
        kwargs.i = "";
      end
      self.i = kwargs.i;
    end

    function res = eq(self, other)
      res = ...
        isa(other, "basic_types.RecordWithString") && ...
        isequal(self.i, other.i);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = basic_types.RecordWithString();
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