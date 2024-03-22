classdef RecordWithVlenCollections < handle
  properties
    vector
    array
  end

  methods
    function obj = RecordWithVlenCollections(vector, array)
      if nargin > 0
        obj.vector = vector;
        obj.array = array;
      else
        obj.vector = int32.empty();
        obj.array = int32.empty(0, 0);
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithVlenCollections') && ...
        all([obj.vector] == [other.vector]) && ...
        isequal(obj.array, other.array);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithVlenCollections();
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
