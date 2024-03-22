classdef RecordWithArraysSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithArraysSerializer()
      field_serializers{1} = yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer);
      field_serializers{2} = yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer);
      field_serializers{3} = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 1);
      field_serializers{4} = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      field_serializers{5} = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      field_serializers{6} = yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [4, 3]);
      field_serializers{7} = yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [4, 3]);
      field_serializers{8} = yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer);
      field_serializers{9} = yardl.binary.FixedNDArraySerializer(yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 4), [5]);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithArrays', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithArrays'));
      obj.write_(outstream, value.default_array, value.default_array_with_empty_dimension, value.rank_1_array, value.rank_2_array, value.rank_2_array_with_named_dimensions, value.rank_2_fixed_array, value.rank_2_fixed_array_with_named_dimensions, value.dynamic_array, value.array_of_vectors)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithArrays(field_values{:});
    end
  end
end
