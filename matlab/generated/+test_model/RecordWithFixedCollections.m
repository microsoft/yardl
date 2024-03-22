classdef RecordWithFixedCollections < handle
  properties
    fixed_vector
    fixed_array
  end

  methods
    function obj = RecordWithFixedCollections(fixed_vector, fixed_array)
      if nargin > 0
        obj.fixed_vector = fixed_vector;
        obj.fixed_array = fixed_array;
      else
        obj.fixed_vector = repelem(int32(0), 3);
        obj.fixed_array = repelem(int32(0), 3, 2);
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithFixedCollections') && ...
        all([obj.fixed_vector] == [other.fixed_vector]) && ...
        isequal(obj.fixed_array, other.fixed_array);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithFixedCollections();
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
