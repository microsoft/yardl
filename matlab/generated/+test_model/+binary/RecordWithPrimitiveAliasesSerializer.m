classdef RecordWithPrimitiveAliasesSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithPrimitiveAliasesSerializer()
      field_serializers{1} = yardl.binary.Uint8Serializer;
      field_serializers{2} = yardl.binary.Int32Serializer;
      field_serializers{3} = yardl.binary.Uint32Serializer;
      field_serializers{4} = yardl.binary.Int64Serializer;
      field_serializers{5} = yardl.binary.Uint64Serializer;
      field_serializers{6} = yardl.binary.Float32Serializer;
      field_serializers{7} = yardl.binary.Float64Serializer;
      field_serializers{8} = yardl.binary.Complexfloat32Serializer;
      field_serializers{9} = yardl.binary.Complexfloat64Serializer;
      obj@yardl.binary.RecordSerializer('test_model.RecordWithPrimitiveAliases', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithPrimitiveAliases'));
      obj.write_(outstream, value.byte_field, value.int_field, value.uint_field, value.long_field, value.ulong_field, value.float_field, value.double_field, value.complexfloat_field, value.complexdouble_field)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithPrimitiveAliases(field_values{:});
    end
  end
end
