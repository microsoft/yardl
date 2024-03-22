classdef RecordWithComputedFields < handle
  properties
    array_field
    array_field_map_dimensions
    dynamic_array_field
    fixed_array_field
    int_field
    int8_field
    uint8_field
    int16_field
    uint16_field
    uint32_field
    int64_field
    uint64_field
    size_field
    float32_field
    float64_field
    complexfloat32_field
    complexfloat64_field
    string_field
    tuple_field
    vector_field
    vector_of_vectors_field
    fixed_vector_field
    optional_named_array
    int_float_union
    nullable_int_float_union
    union_with_nested_generic_union
    map_field
  end

  methods
    function obj = RecordWithComputedFields(array_field, array_field_map_dimensions, dynamic_array_field, fixed_array_field, int_field, int8_field, uint8_field, int16_field, uint16_field, uint32_field, int64_field, uint64_field, size_field, float32_field, float64_field, complexfloat32_field, complexfloat64_field, string_field, tuple_field, vector_field, vector_of_vectors_field, fixed_vector_field, optional_named_array, int_float_union, nullable_int_float_union, union_with_nested_generic_union, map_field)
      if nargin > 0
        obj.array_field = array_field;
        obj.array_field_map_dimensions = array_field_map_dimensions;
        obj.dynamic_array_field = dynamic_array_field;
        obj.fixed_array_field = fixed_array_field;
        obj.int_field = int_field;
        obj.int8_field = int8_field;
        obj.uint8_field = uint8_field;
        obj.int16_field = int16_field;
        obj.uint16_field = uint16_field;
        obj.uint32_field = uint32_field;
        obj.int64_field = int64_field;
        obj.uint64_field = uint64_field;
        obj.size_field = size_field;
        obj.float32_field = float32_field;
        obj.float64_field = float64_field;
        obj.complexfloat32_field = complexfloat32_field;
        obj.complexfloat64_field = complexfloat64_field;
        obj.string_field = string_field;
        obj.tuple_field = tuple_field;
        obj.vector_field = vector_field;
        obj.vector_of_vectors_field = vector_of_vectors_field;
        obj.fixed_vector_field = fixed_vector_field;
        obj.optional_named_array = optional_named_array;
        obj.int_float_union = int_float_union;
        obj.nullable_int_float_union = nullable_int_float_union;
        obj.union_with_nested_generic_union = union_with_nested_generic_union;
        obj.map_field = map_field;
      else
        obj.array_field = int32.empty(0, 0);
        obj.array_field_map_dimensions = int32.empty(0, 0);
        obj.dynamic_array_field = int32.empty();
        obj.fixed_array_field = repelem(int32(0), 4, 3);
        obj.int_field = int32(0);
        obj.int8_field = int8(0);
        obj.uint8_field = uint8(0);
        obj.int16_field = int16(0);
        obj.uint16_field = uint16(0);
        obj.uint32_field = uint32(0);
        obj.int64_field = int64(0);
        obj.uint64_field = uint64(0);
        obj.size_field = uint64(0);
        obj.float32_field = single(0);
        obj.float64_field = double(0);
        obj.complexfloat32_field = complex(single(0));
        obj.complexfloat64_field = complex(0);
        obj.string_field = "";
        obj.tuple_field = tuples.Tuple(int32(0), int32(0));
        obj.vector_field = int32.empty();
        obj.vector_of_vectors_field = int32.empty();
        obj.fixed_vector_field = repelem(int32(0), 3);
        obj.optional_named_array = yardl.None;
        obj.int_float_union = test_model.Int32OrFloat32.Int32(int32(0));
        obj.nullable_int_float_union = yardl.None;
        obj.union_with_nested_generic_union = test_model.IntOrGenericRecordWithComputedFields.Int(int32(0));
        obj.map_field = dictionary;
      end
    end

    function res = int_literal(self)
      res = 42;
      return
    end

    function res = large_negative_int64_literal(self)
      res = -4611686018427387904;
      return
    end

    function res = large_u_int64_literal(self)
      res = 9223372036854775808;
      return
    end

    function res = string_literal(self)
      res = "hello";
      return
    end

    function res = string_literal_2(self)
      res = "hello";
      return
    end

    function res = string_literal_3(self)
      res = "hello";
      return
    end

    function res = string_literal_4(self)
      res = "hello";
      return
    end

    function res = access_other_computed_field(self)
      res = self.int_field;
      return
    end

    function res = access_int_field(self)
      res = self.int_field;
      return
    end

    function res = access_string_field(self)
      res = self.string_field;
      return
    end

    function res = access_tuple_field(self)
      res = self.tuple_field;
      return
    end

    function res = access_nested_tuple_field(self)
      res = self.tuple_field.v2;
      return
    end

    function res = access_array_field(self)
      res = self.array_field;
      return
    end

    function res = access_array_field_element(self)
      res = self.array_field(1+0, 1+1);
      return
    end

    function res = access_array_field_element_by_name(self)
      res = self.array_field(1+0, 1+1);
      return
    end

    function res = access_vector_field(self)
      res = self.vector_field;
      return
    end

    function res = access_vector_field_element(self)
      res = self.vector_field(1+1);
      return
    end

    function res = access_vector_of_vectors_field(self)
      res = self.vector_of_vectors_field(1+1, 1+2);
      return
    end

    function res = array_size(self)
      res = numel(self.array_field);
      return
    end

    function res = array_x_size(self)
      res = size(self.array_field, 1+0);
      return
    end

    function res = array_y_size(self)
      res = size(self.array_field, 1+1);
      return
    end

    function res = array_0_size(self)
      res = size(self.array_field, 1+0);
      return
    end

    function res = array_1_size(self)
      res = size(self.array_field, 1+1);
      return
    end

    function res = array_size_from_int_field(self)
      res = size(self.array_field, 1+self.int_field);
      return
    end

    function res = array_size_from_string_field(self)
      function dim = helper_0_(dim_name)
        if dim_name == "x"
          dim = 0;
        elseif dim_name == "y"
          dim = 1;
        else
          throw(yardl.KeyError("Unknown dimension name: '%s'", dim_name));
        end

      end
      res = size(self.array_field, 1+1 + helper_0_(self.string_field));
      return
    end

    function res = array_size_from_nested_int_field(self)
      res = size(self.array_field, 1+self.tuple_field.v1);
      return
    end

    function res = array_field_map_dimensions_x_size(self)
      res = size(self.array_field_map_dimensions, 1+0);
      return
    end

    function res = fixed_array_size(self)
      res = 12;
      return
    end

    function res = fixed_array_x_size(self)
      res = 3;
      return
    end

    function res = fixed_array_0_size(self)
      res = 3;
      return
    end

    function res = vector_size(self)
      res = length(self.vector_field);
      return
    end

    function res = fixed_vector_size(self)
      res = 3;
      return
    end

    function res = array_dimension_x_index(self)
      function dim = helper_0_(dim_name)
        if dim_name == "x"
          dim = 0;
        elseif dim_name == "y"
          dim = 1;
        else
          throw(yardl.KeyError("Unknown dimension name: '%s'", dim_name));
        end

      end
      res = 1 + helper_0_("x");
      return
    end

    function res = array_dimension_y_index(self)
      function dim = helper_0_(dim_name)
        if dim_name == "x"
          dim = 0;
        elseif dim_name == "y"
          dim = 1;
        else
          throw(yardl.KeyError("Unknown dimension name: '%s'", dim_name));
        end

      end
      res = 1 + helper_0_("y");
      return
    end

    function res = array_dimension_index_from_string_field(self)
      function dim = helper_0_(dim_name)
        if dim_name == "x"
          dim = 0;
        elseif dim_name == "y"
          dim = 1;
        else
          throw(yardl.KeyError("Unknown dimension name: '%s'", dim_name));
        end

      end
      res = 1 + helper_0_(self.string_field);
      return
    end

    function res = array_dimension_count(self)
      res = 2;
      return
    end

    function res = dynamic_array_dimension_count(self)
      res = yardl.dimension_count(self.dynamic_array_field);
      return
    end

    function res = access_map(self)
      res = self.map_field;
      return
    end

    function res = map_size(self)
      res = numEntries(self.map_field);
      return
    end

    function res = access_map_entry(self)
      res = self.map_field("hello");
      return
    end

    function res = string_computed_field(self)
      res = "hello";
      return
    end

    function res = access_map_entry_with_computed_field(self)
      res = self.map_field(self.string_computed_field());
      return
    end

    function res = access_map_entry_with_computed_field_nested(self)
      res = self.map_field(self.map_field(self.string_computed_field()));
      return
    end

    function res = access_missing_map_entry(self)
      res = self.map_field("missing");
      return
    end

    function res = optional_named_array_length(self)
      var1 = self.optional_named_array;
      if var1 ~= yardl.None
        arr = var1;
        res = numel(arr);
        return
      end
      res = 0;
      return
    end

    function res = optional_named_array_length_with_discard(self)
      var1 = self.optional_named_array;
      if var1 ~= yardl.None
        arr = var1;
        res = numel(arr);
        return
      end
      res = 0;
      return
    end

    function res = int_float_union_as_float(self)
      var1 = self.int_float_union;
      if var1.index == 1
        i_foo = var1.value;
        res = single(i_foo);
        return
      end
      if var1.index == 2
        f = var1.value;
        res = f;
        return
      end
      throw(yardl.RuntimeError("Unexpected union case"))
    end

    function res = nullable_int_float_union_string(self)
      var1 = self.nullable_int_float_union;
      if var1 == yardl.None
        res = "null";
        return
      end
      if var1.index == 1
        res = "int";
        return
      end
      res = "float";
      return
      throw(yardl.RuntimeError("Unexpected union case"))
    end

    function res = nested_switch(self)
      var1 = self.union_with_nested_generic_union;
      if var1.index == 1
        res = -1;
        return
      end
      if var1.index == 2
        rec = var1.value;
        var2 = rec.f1;
        if var2.index == 2
          res = int32(20);
          return
        end
        if var2.index == 1
          res = int32(10);
          return
        end
        throw(yardl.RuntimeError("Unexpected union case"))
      end
      throw(yardl.RuntimeError("Unexpected union case"))
    end

    function res = use_nested_computed_field(self)
      var1 = self.union_with_nested_generic_union;
      if var1.index == 1
        res = -1;
        return
      end
      if var1.index == 2
        rec = var1.value;
        res = int32(rec.type_index());
        return
      end
      throw(yardl.RuntimeError("Unexpected union case"))
    end

    function res = switch_over_single_value(self)
      var1 = self.int_field;
      i = var1;
      res = i;
      return
    end

    function res = arithmetic_1(self)
      res = 1 + 2;
      return
    end

    function res = arithmetic_2(self)
      res = 1 + 2 .* 3 + 4;
      return
    end

    function res = arithmetic_3(self)
      res = (1 + 2) .* 3 + 4;
      return
    end

    function res = arithmetic_4(self)
      res = self.array_size_from_int_field() + 2;
      return
    end

    function res = arithmetic_5(self)
      res = size(self.array_field, 1+2 - 1);
      return
    end

    function res = arithmetic_6(self)
      res = 7 ./ 2;
      return
    end

    function res = arithmetic_7(self)
      res = double(7) ^ double(2);
      return
    end

    function res = arithmetic8(self)
      res = self.complexfloat32_field .* complex(single(single(3)));
      return
    end

    function res = arithmetic_9(self)
      res = 1.2 + double(1);
      return
    end

    function res = arithmetic_10(self)
      res = 1e10 + 9e9;
      return
    end

    function res = arithmetic_11(self)
      res = -(4.3 + double(1));
      return
    end

    function res = cast_int_to_float(self)
      res = single(self.int_field);
      return
    end

    function res = cast_float_to_int(self)
      res = int32(self.float32_field);
      return
    end

    function res = cast_power(self)
      res = int32(double(7) ^ double(2));
      return
    end

    function res = cast_complex32_to_complex64(self)
      res = complex(double(self.complexfloat32_field));
      return
    end

    function res = cast_complex64_to_complex32(self)
      res = complex(single(self.complexfloat64_field));
      return
    end

    function res = cast_float_to_complex(self)
      res = complex(single(66.6));
      return
    end


    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithComputedFields') && ...
        isequal(obj.array_field, other.array_field) && ...
        isequal(obj.array_field_map_dimensions, other.array_field_map_dimensions) && ...
        isequal(obj.dynamic_array_field, other.dynamic_array_field) && ...
        isequal(obj.fixed_array_field, other.fixed_array_field) && ...
        all([obj.int_field] == [other.int_field]) && ...
        all([obj.int8_field] == [other.int8_field]) && ...
        all([obj.uint8_field] == [other.uint8_field]) && ...
        all([obj.int16_field] == [other.int16_field]) && ...
        all([obj.uint16_field] == [other.uint16_field]) && ...
        all([obj.uint32_field] == [other.uint32_field]) && ...
        all([obj.int64_field] == [other.int64_field]) && ...
        all([obj.uint64_field] == [other.uint64_field]) && ...
        all([obj.size_field] == [other.size_field]) && ...
        all([obj.float32_field] == [other.float32_field]) && ...
        all([obj.float64_field] == [other.float64_field]) && ...
        all([obj.complexfloat32_field] == [other.complexfloat32_field]) && ...
        all([obj.complexfloat64_field] == [other.complexfloat64_field]) && ...
        all([obj.string_field] == [other.string_field]) && ...
        isequal(obj.tuple_field, other.tuple_field) && ...
        all([obj.vector_field] == [other.vector_field]) && ...
        all([obj.vector_of_vectors_field] == [other.vector_of_vectors_field]) && ...
        all([obj.fixed_vector_field] == [other.fixed_vector_field]) && ...
        isequal(obj.optional_named_array, other.optional_named_array) && ...
        all([obj.int_float_union] == [other.int_float_union]) && ...
        all([obj.nullable_int_float_union] == [other.nullable_int_float_union]) && ...
        all([obj.union_with_nested_generic_union] == [other.union_with_nested_generic_union]) && ...
        all([obj.map_field] == [other.map_field]);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithComputedFields();
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
