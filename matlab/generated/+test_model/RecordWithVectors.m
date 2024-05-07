% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithVectors < handle
  properties
    default_vector
    default_vector_fixed_length
    vector_of_vectors
  end

  methods
    function self = RecordWithVectors(kwargs)
      arguments
        kwargs.default_vector = int32.empty();
        kwargs.default_vector_fixed_length = repelem(int32(0), 3);
        kwargs.vector_of_vectors = int32.empty();
      end
      self.default_vector = kwargs.default_vector;
      self.default_vector_fixed_length = kwargs.default_vector_fixed_length;
      self.vector_of_vectors = kwargs.vector_of_vectors;
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordWithVectors") && ...
        all([self.default_vector] == [other.default_vector]) && ...
        all([self.default_vector_fixed_length] == [other.default_vector_fixed_length]) && ...
        all([self.vector_of_vectors] == [other.vector_of_vectors]);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithVectors();
      if nargin == 0
        z = elem;
        return;
      end
      sz = [varargin{:}];
      if isscalar(sz)
        sz = [sz, sz];
      end
      z = reshape(repelem(elem, prod(sz)), sz);
    end
  end
end
