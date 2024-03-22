classdef SimpleRecord < handle
  properties
    x
    y
    z
  end

  methods
    function obj = SimpleRecord(x, y, z)
      if nargin > 0
        obj.x = x;
        obj.y = y;
        obj.z = z;
      else
        obj.x = int32(0);
        obj.y = int32(0);
        obj.z = int32(0);
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.SimpleRecord') && ...
        all([obj.x] == [other.x]) && ...
        all([obj.y] == [other.y]) && ...
        all([obj.z] == [other.z]);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.SimpleRecord();
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