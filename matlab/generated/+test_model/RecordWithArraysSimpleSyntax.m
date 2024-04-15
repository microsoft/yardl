% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithArraysSimpleSyntax < handle
  properties
    default_array
    default_array_with_empty_dimension
    rank_1_array
    rank_2_array
    rank_2_array_with_named_dimensions
    rank_2_fixed_array
    rank_2_fixed_array_with_named_dimensions
    dynamic_array
    array_of_vectors
  end

  methods
    function self = RecordWithArraysSimpleSyntax(default_array, default_array_with_empty_dimension, rank_1_array, rank_2_array, rank_2_array_with_named_dimensions, rank_2_fixed_array, rank_2_fixed_array_with_named_dimensions, dynamic_array, array_of_vectors)
      if nargin > 0
        self.default_array = default_array;
        self.default_array_with_empty_dimension = default_array_with_empty_dimension;
        self.rank_1_array = rank_1_array;
        self.rank_2_array = rank_2_array;
        self.rank_2_array_with_named_dimensions = rank_2_array_with_named_dimensions;
        self.rank_2_fixed_array = rank_2_fixed_array;
        self.rank_2_fixed_array_with_named_dimensions = rank_2_fixed_array_with_named_dimensions;
        self.dynamic_array = dynamic_array;
        self.array_of_vectors = array_of_vectors;
      else
        self.default_array = int32.empty();
        self.default_array_with_empty_dimension = int32.empty();
        self.rank_1_array = int32.empty(0);
        self.rank_2_array = int32.empty(0, 0);
        self.rank_2_array_with_named_dimensions = int32.empty(0, 0);
        self.rank_2_fixed_array = repelem(int32(0), 4, 3);
        self.rank_2_fixed_array_with_named_dimensions = repelem(int32(0), 4, 3);
        self.dynamic_array = int32.empty();
        self.array_of_vectors = repelem(repelem(int32(0), 4), 5, 1);
      end
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordWithArraysSimpleSyntax") && ...
        isequal(self.default_array, other.default_array) && ...
        isequal(self.default_array_with_empty_dimension, other.default_array_with_empty_dimension) && ...
        isequal(self.rank_1_array, other.rank_1_array) && ...
        isequal(self.rank_2_array, other.rank_2_array) && ...
        isequal(self.rank_2_array_with_named_dimensions, other.rank_2_array_with_named_dimensions) && ...
        isequal(self.rank_2_fixed_array, other.rank_2_fixed_array) && ...
        isequal(self.rank_2_fixed_array_with_named_dimensions, other.rank_2_fixed_array_with_named_dimensions) && ...
        isequal(self.dynamic_array, other.dynamic_array) && ...
        isequal(self.array_of_vectors, other.array_of_vectors);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithArraysSimpleSyntax();
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
