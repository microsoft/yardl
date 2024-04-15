% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithOptionalVector < handle
  properties
    optional_vector
  end

  methods
    function self = RecordWithOptionalVector(optional_vector)
      if nargin > 0
        self.optional_vector = optional_vector;
      else
        self.optional_vector = yardl.None;
      end
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordWithOptionalVector") && ...
        all([self.optional_vector] == [other.optional_vector]);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithOptionalVector();
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
