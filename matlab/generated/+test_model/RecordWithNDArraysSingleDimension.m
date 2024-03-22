classdef RecordWithNDArraysSingleDimension < handle
  properties
    ints
    fixed_simple_record_array
    fixed_record_with_vlens_array
  end

  methods
    function obj = RecordWithNDArraysSingleDimension(ints, fixed_simple_record_array, fixed_record_with_vlens_array)
      if nargin > 0
        obj.ints = ints;
        obj.fixed_simple_record_array = fixed_simple_record_array;
        obj.fixed_record_with_vlens_array = fixed_record_with_vlens_array;
      else
        obj.ints = int32.empty(0);
        obj.fixed_simple_record_array = test_model.SimpleRecord.empty(0);
        obj.fixed_record_with_vlens_array = test_model.RecordWithVlens.empty(0);
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithNDArraysSingleDimension') && ...
        isequal(obj.ints, other.ints) && ...
        isequal(obj.fixed_simple_record_array, other.fixed_simple_record_array) && ...
        isequal(obj.fixed_record_with_vlens_array, other.fixed_record_with_vlens_array);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithNDArraysSingleDimension();
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
