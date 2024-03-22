classdef RecordWithFixedVectors < handle
  properties
    fixed_int_vector
    fixed_simple_record_vector
    fixed_record_with_vlens_vector
  end

  methods
    function obj = RecordWithFixedVectors(fixed_int_vector, fixed_simple_record_vector, fixed_record_with_vlens_vector)
      if nargin > 0
        obj.fixed_int_vector = fixed_int_vector;
        obj.fixed_simple_record_vector = fixed_simple_record_vector;
        obj.fixed_record_with_vlens_vector = fixed_record_with_vlens_vector;
      else
        obj.fixed_int_vector = repelem(int32(0), 5);
        obj.fixed_simple_record_vector = repelem(test_model.SimpleRecord(), 3);
        obj.fixed_record_with_vlens_vector = repelem(test_model.RecordWithVlens(), 2);
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithFixedVectors') && ...
        all([obj.fixed_int_vector] == [other.fixed_int_vector]) && ...
        all([obj.fixed_simple_record_vector] == [other.fixed_simple_record_vector]) && ...
        all([obj.fixed_record_with_vlens_vector] == [other.fixed_record_with_vlens_vector]);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithFixedVectors();
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
