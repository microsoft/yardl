classdef RecordWithGenericMaps < handle
  properties
    m
    am
  end

  methods
    function obj = RecordWithGenericMaps(m, am)
      if nargin > 0
        obj.m = m;
        obj.am = am;
      else
        obj.m = dictionary;
        obj.am = dictionary;
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithGenericMaps') && ...
        isequal(obj.m, other.m) && ...
        isequal(obj.am, other.am);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithGenericMaps();
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
