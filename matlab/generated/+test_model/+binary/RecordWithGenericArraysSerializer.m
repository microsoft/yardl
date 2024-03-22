classdef RecordWithGenericArraysSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithGenericArraysSerializer(t_serializer)
      field_serializers{1} = yardl.binary.NDArraySerializer(t_serializer, 2);
      field_serializers{2} = yardl.binary.FixedNDArraySerializer(t_serializer, [8, 16]);
      field_serializers{3} = yardl.binary.DynamicNDArraySerializer(t_serializer);
      field_serializers{4} = yardl.binary.NDArraySerializer(t_serializer, 2);
      field_serializers{5} = yardl.binary.FixedNDArraySerializer(t_serializer, [8, 16]);
      field_serializers{6} = yardl.binary.DynamicNDArraySerializer(t_serializer);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithGenericArrays', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithGenericArrays'));
      obj.write_(outstream, value.nd, value.fixed_nd, value.dynamic_nd, value.aliased_nd, value.aliased_fixed_nd, value.aliased_dynamic_nd)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithGenericArrays(field_values{:});
    end
  end
end
