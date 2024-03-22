classdef RecordWithPrimitives < handle
  properties
    bool_field
    int8_field
    uint8_field
    int16_field
    uint16_field
    int32_field
    uint32_field
    int64_field
    uint64_field
    size_field
    float32_field
    float64_field
    complexfloat32_field
    complexfloat64_field
    date_field
    time_field
    datetime_field
  end

  methods
    function obj = RecordWithPrimitives(bool_field, int8_field, uint8_field, int16_field, uint16_field, int32_field, uint32_field, int64_field, uint64_field, size_field, float32_field, float64_field, complexfloat32_field, complexfloat64_field, date_field, time_field, datetime_field)
      if nargin > 0
        obj.bool_field = bool_field;
        obj.int8_field = int8_field;
        obj.uint8_field = uint8_field;
        obj.int16_field = int16_field;
        obj.uint16_field = uint16_field;
        obj.int32_field = int32_field;
        obj.uint32_field = uint32_field;
        obj.int64_field = int64_field;
        obj.uint64_field = uint64_field;
        obj.size_field = size_field;
        obj.float32_field = float32_field;
        obj.float64_field = float64_field;
        obj.complexfloat32_field = complexfloat32_field;
        obj.complexfloat64_field = complexfloat64_field;
        obj.date_field = date_field;
        obj.time_field = time_field;
        obj.datetime_field = datetime_field;
      else
        obj.bool_field = false;
        obj.int8_field = int8(0);
        obj.uint8_field = uint8(0);
        obj.int16_field = int16(0);
        obj.uint16_field = uint16(0);
        obj.int32_field = int32(0);
        obj.uint32_field = uint32(0);
        obj.int64_field = int64(0);
        obj.uint64_field = uint64(0);
        obj.size_field = uint64(0);
        obj.float32_field = single(0);
        obj.float64_field = double(0);
        obj.complexfloat32_field = complex(single(0));
        obj.complexfloat64_field = complex(0);
        obj.date_field = yardl.Date();
        obj.time_field = yardl.Time();
        obj.datetime_field = yardl.DateTime();
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithPrimitives') && ...
        all([obj.bool_field] == [other.bool_field]) && ...
        all([obj.int8_field] == [other.int8_field]) && ...
        all([obj.uint8_field] == [other.uint8_field]) && ...
        all([obj.int16_field] == [other.int16_field]) && ...
        all([obj.uint16_field] == [other.uint16_field]) && ...
        all([obj.int32_field] == [other.int32_field]) && ...
        all([obj.uint32_field] == [other.uint32_field]) && ...
        all([obj.int64_field] == [other.int64_field]) && ...
        all([obj.uint64_field] == [other.uint64_field]) && ...
        all([obj.size_field] == [other.size_field]) && ...
        all([obj.float32_field] == [other.float32_field]) && ...
        all([obj.float64_field] == [other.float64_field]) && ...
        all([obj.complexfloat32_field] == [other.complexfloat32_field]) && ...
        all([obj.complexfloat64_field] == [other.complexfloat64_field]) && ...
        all([obj.date_field] == [other.date_field]) && ...
        all([obj.time_field] == [other.time_field]) && ...
        all([obj.datetime_field] == [other.datetime_field]);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithPrimitives();
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
