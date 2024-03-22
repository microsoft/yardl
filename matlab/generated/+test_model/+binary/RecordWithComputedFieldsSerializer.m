classdef RecordWithComputedFieldsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithComputedFieldsSerializer()
      field_serializers{1} = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      field_serializers{2} = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      field_serializers{3} = yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer);
      field_serializers{4} = yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [4, 3]);
      field_serializers{5} = yardl.binary.Int32Serializer;
      field_serializers{6} = yardl.binary.Int8Serializer;
      field_serializers{7} = yardl.binary.Uint8Serializer;
      field_serializers{8} = yardl.binary.Int16Serializer;
      field_serializers{9} = yardl.binary.Uint16Serializer;
      field_serializers{10} = yardl.binary.Uint32Serializer;
      field_serializers{11} = yardl.binary.Int64Serializer;
      field_serializers{12} = yardl.binary.Uint64Serializer;
      field_serializers{13} = yardl.binary.SizeSerializer;
      field_serializers{14} = yardl.binary.Float32Serializer;
      field_serializers{15} = yardl.binary.Float64Serializer;
      field_serializers{16} = yardl.binary.Complexfloat32Serializer;
      field_serializers{17} = yardl.binary.Complexfloat64Serializer;
      field_serializers{18} = yardl.binary.StringSerializer;
      field_serializers{19} = tuples.binary.TupleSerializer(yardl.binary.Int32Serializer, yardl.binary.Int32Serializer);
      field_serializers{20} = yardl.binary.VectorSerializer(yardl.binary.Int32Serializer);
      field_serializers{21} = yardl.binary.VectorSerializer(yardl.binary.VectorSerializer(yardl.binary.Int32Serializer));
      field_serializers{22} = yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 3);
      field_serializers{23} = yardl.binary.OptionalSerializer(yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2));
      field_serializers{24} = yardl.binary.UnionSerializer('test_model.Int32OrFloat32', {yardl.binary.Int32Serializer, yardl.binary.Float32Serializer}, {@test_model.Int32OrFloat32.Int32, @test_model.Int32OrFloat32.Float32});
      field_serializers{25} = yardl.binary.UnionSerializer('test_model.Int32OrFloat32', {yardl.binary.NoneSerializer, yardl.binary.Int32Serializer, yardl.binary.Float32Serializer}, {yardl.None, @test_model.Int32OrFloat32.Int32, @test_model.Int32OrFloat32.Float32});
      field_serializers{26} = yardl.binary.UnionSerializer('test_model.IntOrGenericRecordWithComputedFields', {yardl.binary.Int32Serializer, basic_types.binary.GenericRecordWithComputedFieldsSerializer(yardl.binary.StringSerializer, yardl.binary.Float32Serializer)}, {@test_model.IntOrGenericRecordWithComputedFields.Int, @test_model.IntOrGenericRecordWithComputedFields.GenericRecordWithComputedFields});
      field_serializers{27} = yardl.binary.MapSerializer(yardl.binary.StringSerializer, yardl.binary.StringSerializer);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithComputedFields', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithComputedFields'));
      obj.write_(outstream, value.array_field, value.array_field_map_dimensions, value.dynamic_array_field, value.fixed_array_field, value.int_field, value.int8_field, value.uint8_field, value.int16_field, value.uint16_field, value.uint32_field, value.int64_field, value.uint64_field, value.size_field, value.float32_field, value.float64_field, value.complexfloat32_field, value.complexfloat64_field, value.string_field, value.tuple_field, value.vector_field, value.vector_of_vectors_field, value.fixed_vector_field, value.optional_named_array, value.int_float_union, value.nullable_int_float_union, value.union_with_nested_generic_union, value.map_field)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithComputedFields(field_values{:});
    end
  end
end