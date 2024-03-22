classdef RecordWithUnionsOfContainers < handle
  properties
    map_or_scalar
    vector_or_scalar
    array_or_scalar
  end

  methods
    function obj = RecordWithUnionsOfContainers(map_or_scalar, vector_or_scalar, array_or_scalar)
      if nargin > 0
        obj.map_or_scalar = map_or_scalar;
        obj.vector_or_scalar = vector_or_scalar;
        obj.array_or_scalar = array_or_scalar;
      else
        obj.map_or_scalar = test_model.MapOrScalar.Map(dictionary);
        obj.vector_or_scalar = test_model.VectorOrScalar.Vector(int32.empty());
        obj.array_or_scalar = test_model.ArrayOrScalar.Array(int32.empty());
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithUnionsOfContainers') && ...
        all([obj.map_or_scalar] == [other.map_or_scalar]) && ...
        all([obj.vector_or_scalar] == [other.vector_or_scalar]) && ...
        isequal(obj.array_or_scalar, other.array_or_scalar);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithUnionsOfContainers();
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
