classdef RecordWithVectors < handle
  properties
    default_vector
    default_vector_fixed_length
    vector_of_vectors
  end

  methods
    function obj = RecordWithVectors(default_vector, default_vector_fixed_length, vector_of_vectors)
      if nargin > 0
        obj.default_vector = default_vector;
        obj.default_vector_fixed_length = default_vector_fixed_length;
        obj.vector_of_vectors = vector_of_vectors;
      else
        obj.default_vector = int32.empty();
        obj.default_vector_fixed_length = repelem(int32(0), 3);
        obj.vector_of_vectors = int32.empty();
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithVectors') && ...
        all([obj.default_vector] == [other.default_vector]) && ...
        all([obj.default_vector_fixed_length] == [other.default_vector_fixed_length]) && ...
        all([obj.vector_of_vectors] == [other.vector_of_vectors]);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithVectors();
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
