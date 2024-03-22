classdef RecordWithAliasedOptionalGenericField < handle
  properties
    v
  end

  methods
    function obj = RecordWithAliasedOptionalGenericField(v)
      if nargin > 0
        obj.v = v;
      else
        obj.v = yardl.None;
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithAliasedOptionalGenericField') && ...
        isequal(obj.v, other.v);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithAliasedOptionalGenericField();
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
