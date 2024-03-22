classdef SimpleEncodingCounters < handle
  properties
    e1
    e2
    slice
    repetition
  end

  methods
    function obj = SimpleEncodingCounters(e1, e2, slice, repetition)
      if nargin > 0
        obj.e1 = e1;
        obj.e2 = e2;
        obj.slice = slice;
        obj.repetition = repetition;
      else
        obj.e1 = yardl.None;
        obj.e2 = yardl.None;
        obj.slice = yardl.None;
        obj.repetition = yardl.None;
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.SimpleEncodingCounters') && ...
        all([obj.e1] == [other.e1]) && ...
        all([obj.e2] == [other.e2]) && ...
        all([obj.slice] == [other.slice]) && ...
        all([obj.repetition] == [other.repetition]);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.SimpleEncodingCounters();
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
