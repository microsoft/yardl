% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithNamedFixedArrays < handle
  properties
    ints
    fixed_simple_record_array
    fixed_record_with_vlens_array
  end

  methods
    function self = RecordWithNamedFixedArrays(ints, fixed_simple_record_array, fixed_record_with_vlens_array)
      if nargin > 0
        self.ints = ints;
        self.fixed_simple_record_array = fixed_simple_record_array;
        self.fixed_record_with_vlens_array = fixed_record_with_vlens_array;
      else
        self.ints = repelem(int32(0), 3, 2);
        self.fixed_simple_record_array = repelem(test_model.SimpleRecord(), 2, 3);
        self.fixed_record_with_vlens_array = repelem(test_model.RecordWithVlens(), 2, 2);
      end
    end

    function res = eq(self, other)
      res = ...
        isa(other, 'test_model.RecordWithNamedFixedArrays') && ...
        isequal(self.ints, other.ints) && ...
        isequal(self.fixed_simple_record_array, other.fixed_simple_record_array) && ...
        isequal(self.fixed_record_with_vlens_array, other.fixed_record_with_vlens_array);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithNamedFixedArrays();
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
