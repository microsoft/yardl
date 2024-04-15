% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithFixedVectors < handle
  properties
    fixed_int_vector
    fixed_simple_record_vector
    fixed_record_with_vlens_vector
  end

  methods
    function self = RecordWithFixedVectors(fixed_int_vector, fixed_simple_record_vector, fixed_record_with_vlens_vector)
      if nargin > 0
        self.fixed_int_vector = fixed_int_vector;
        self.fixed_simple_record_vector = fixed_simple_record_vector;
        self.fixed_record_with_vlens_vector = fixed_record_with_vlens_vector;
      else
        self.fixed_int_vector = repelem(int32(0), 5);
        self.fixed_simple_record_vector = repelem(test_model.SimpleRecord(), 3);
        self.fixed_record_with_vlens_vector = repelem(test_model.RecordWithVlens(), 2);
      end
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordWithFixedVectors") && ...
        all([self.fixed_int_vector] == [other.fixed_int_vector]) && ...
        all([self.fixed_simple_record_vector] == [other.fixed_simple_record_vector]) && ...
        all([self.fixed_record_with_vlens_vector] == [other.fixed_record_with_vlens_vector]);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithFixedVectors();
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
