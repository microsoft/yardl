% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithGenericFixedVectors < handle
  properties
    fv
    afv
  end

  methods
    function self = RecordWithGenericFixedVectors(kwargs)
      arguments
        kwargs.fv;
        kwargs.afv;
      end
      if ~isfield(kwargs, "fv")
        throw(yardl.TypeError("Missing required keyword argument 'fv'"))
      end
      self.fv = kwargs.fv;
      if ~isfield(kwargs, "afv")
        throw(yardl.TypeError("Missing required keyword argument 'afv'"))
      end
      self.afv = kwargs.afv;
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordWithGenericFixedVectors") && ...
        isequal({self.fv}, {other.fv}) && ...
        isequal({self.afv}, {other.afv});
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
      elem = test_model.RecordWithGenericFixedVectors(fv=yardl.None, afv=yardl.None);
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
