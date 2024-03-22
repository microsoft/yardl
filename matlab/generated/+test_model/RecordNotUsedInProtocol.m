classdef RecordNotUsedInProtocol < handle
  properties
    u1
    u2
  end

  methods
    function obj = RecordNotUsedInProtocol(u1, u2)
      if nargin > 0
        obj.u1 = u1;
        obj.u2 = u2;
      else
        obj.u1 = test_model.GenericUnion3.T(int32(0));
        obj.u2 = test_model.GenericUnion3Alternate.U(int32(0));
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordNotUsedInProtocol') && ...
        all([obj.u1] == [other.u1]) && ...
        all([obj.u2] == [other.u2]);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordNotUsedInProtocol();
      if nargin == 0
        z = elem;
      elseif nargin == 1
        n = varargin{1};
        z = reshape(repelem(elem, n*n), [n, n]);
      else
        sz = [varargin{:}];
        z = reshape(repelem(elem, prod(sz)), sz);
      end
    end
  end
end
