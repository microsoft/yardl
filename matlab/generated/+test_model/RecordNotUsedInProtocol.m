% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordNotUsedInProtocol < handle
  properties
    u1
    u2
  end

  methods
    function self = RecordNotUsedInProtocol(u1, u2)
      if nargin > 0
        self.u1 = u1;
        self.u2 = u2;
      else
        self.u1 = test_model.GenericUnion3.T(int32(0));
        self.u2 = test_model.GenericUnion3Alternate.U(int32(0));
      end
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordNotUsedInProtocol") && ...
        all([self.u1] == [other.u1]) && ...
        all([self.u2] == [other.u2]);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordNotUsedInProtocol();
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
