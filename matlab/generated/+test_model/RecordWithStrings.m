classdef RecordWithStrings < handle
  properties
    a
    b
  end

  methods
    function obj = RecordWithStrings(a, b)
      if nargin > 0
        obj.a = a;
        obj.b = b;
      else
        obj.a = "";
        obj.b = "";
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithStrings') && ...
        all([obj.a] == [other.a]) && ...
        all([obj.b] == [other.b]);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithStrings();
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
