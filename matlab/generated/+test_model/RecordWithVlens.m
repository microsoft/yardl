classdef RecordWithVlens < handle
  properties
    a
    b
    c
  end

  methods
    function obj = RecordWithVlens(a, b, c)
      if nargin > 0
        obj.a = a;
        obj.b = b;
        obj.c = c;
      else
        obj.a = test_model.SimpleRecord.empty();
        obj.b = int32(0);
        obj.c = int32(0);
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithVlens') && ...
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
      elem = test_model.RecordWithVlens();
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