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
    function obj = RecordWithArraysSimpleSyntax(default_array, default_array_with_empty_dimension, rank_1_array, rank_2_array, rank_2_array_with_named_dimensions, rank_2_fixed_array, rank_2_fixed_array_with_named_dimensions, dynamic_array, array_of_vectors)
      if nargin > 0
        obj.default_array = default_array;
        obj.default_array_with_empty_dimension = default_array_with_empty_dimension;
        obj.rank_1_array = rank_1_array;
        obj.rank_2_array = rank_2_array;
        obj.rank_2_array_with_named_dimensions = rank_2_array_with_named_dimensions;
        obj.rank_2_fixed_array = rank_2_fixed_array;
        obj.rank_2_fixed_array_with_named_dimensions = rank_2_fixed_array_with_named_dimensions;
        obj.dynamic_array = dynamic_array;
        obj.array_of_vectors = array_of_vectors;
      else
        obj.default_array = int32.empty();
        obj.default_array_with_empty_dimension = int32.empty();
        obj.rank_1_array = int32.empty(0);
        obj.rank_2_array = int32.empty(0, 0);
        obj.rank_2_array_with_named_dimensions = int32.empty(0, 0);
        obj.rank_2_fixed_array = repelem(int32(0), 4, 3);
        obj.rank_2_fixed_array_with_named_dimensions = repelem(int32(0), 4, 3);
        obj.dynamic_array = int32.empty();
        obj.array_of_vectors = repelem(repelem(int32(0), 4), 5, 1);
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithArraysSimpleSyntax') && ...
        isequal(obj.default_array, other.default_array) && ...
        isequal(obj.default_array_with_empty_dimension, other.default_array_with_empty_dimension) && ...
        isequal(obj.rank_1_array, other.rank_1_array) && ...
        isequal(obj.rank_2_array, other.rank_2_array) && ...
        isequal(obj.rank_2_array_with_named_dimensions, other.rank_2_array_with_named_dimensions) && ...
        isequal(obj.rank_2_fixed_array, other.rank_2_fixed_array) && ...
        isequal(obj.rank_2_fixed_array_with_named_dimensions, other.rank_2_fixed_array_with_named_dimensions) && ...
        isequal(obj.dynamic_array, other.dynamic_array) && ...
        isequal(obj.array_of_vectors, other.array_of_vectors);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithArraysSimpleSyntax();
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
