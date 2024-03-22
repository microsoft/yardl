classdef RecordWithPrimitiveAliases < handle
  properties
    byte_field
    int_field
    uint_field
    long_field
    ulong_field
    float_field
    double_field
    complexfloat_field
    complexdouble_field
  end

  methods
    function obj = RecordWithPrimitiveAliases(byte_field, int_field, uint_field, long_field, ulong_field, float_field, double_field, complexfloat_field, complexdouble_field)
      if nargin > 0
        obj.byte_field = byte_field;
        obj.int_field = int_field;
        obj.uint_field = uint_field;
        obj.long_field = long_field;
        obj.ulong_field = ulong_field;
        obj.float_field = float_field;
        obj.double_field = double_field;
        obj.complexfloat_field = complexfloat_field;
        obj.complexdouble_field = complexdouble_field;
      else
        obj.byte_field = uint8(0);
        obj.int_field = int32(0);
        obj.uint_field = uint32(0);
        obj.long_field = int64(0);
        obj.ulong_field = uint64(0);
        obj.float_field = single(0);
        obj.double_field = double(0);
        obj.complexfloat_field = complex(single(0));
        obj.complexdouble_field = complex(0);
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithPrimitiveAliases') && ...
        all([obj.byte_field] == [other.byte_field]) && ...
        all([obj.int_field] == [other.int_field]) && ...
        all([obj.uint_field] == [other.uint_field]) && ...
        all([obj.long_field] == [other.long_field]) && ...
        all([obj.ulong_field] == [other.ulong_field]) && ...
        all([obj.float_field] == [other.float_field]) && ...
        all([obj.double_field] == [other.double_field]) && ...
        all([obj.complexfloat_field] == [other.complexfloat_field]) && ...
        all([obj.complexdouble_field] == [other.complexdouble_field]);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithPrimitiveAliases();
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