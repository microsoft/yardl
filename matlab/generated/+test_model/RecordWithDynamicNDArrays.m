classdef RecordWithDynamicNDArrays < handle
  properties
    ints
    simple_record_array
    record_with_vlens_array
  end

  methods
    function obj = RecordWithDynamicNDArrays(ints, simple_record_array, record_with_vlens_array)
      if nargin > 0
        obj.ints = ints;
        obj.simple_record_array = simple_record_array;
        obj.record_with_vlens_array = record_with_vlens_array;
      else
        obj.ints = int32.empty();
        obj.simple_record_array = test_model.SimpleRecord.empty();
        obj.record_with_vlens_array = test_model.RecordWithVlens.empty();
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithDynamicNDArrays') && ...
        isequal(obj.ints, other.ints) && ...
        isequal(obj.simple_record_array, other.simple_record_array) && ...
        isequal(obj.record_with_vlens_array, other.record_with_vlens_array);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithDynamicNDArrays();
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
