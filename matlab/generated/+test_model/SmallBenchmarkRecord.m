classdef SmallBenchmarkRecord < handle
  properties
    a
    b
    c
  end

  methods
    function obj = SmallBenchmarkRecord(a, b, c)
      if nargin > 0
        obj.a = a;
        obj.b = b;
        obj.c = c;
      else
        obj.a = double(0);
        obj.b = single(0);
        obj.c = single(0);
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.SmallBenchmarkRecord') && ...
        all([obj.a] == [other.a]) && ...
        all([obj.b] == [other.b]) && ...
        all([obj.c] == [other.c]);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.SmallBenchmarkRecord();
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
