% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithArraysSerializer < yardl.binary.RecordSerializer
  methods
    function self = RecordWithArraysSerializer()
      field_serializers{1} = yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer);
      field_serializers{2} = yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer);
      field_serializers{3} = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 1);
      field_serializers{4} = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      field_serializers{5} = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      field_serializers{6} = yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [4, 3]);
      field_serializers{7} = yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [4, 3]);
      field_serializers{8} = yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer);
      field_serializers{9} = yardl.binary.FixedNDArraySerializer(yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 4), [5]);
      self@yardl.binary.RecordSerializer('test_model.RecordWithArrays', field_serializers);
    end

    function write(self, outstream, value)
      arguments
        self
        outstream (1,1) yardl.binary.CodedOutputStream
        value (1,1) test_model.RecordWithArrays
      end
      self.write_(outstream, value.default_array, value.default_array_with_empty_dimension, value.rank_1_array, value.rank_2_array, value.rank_2_array_with_named_dimensions, value.rank_2_fixed_array, value.rank_2_fixed_array_with_named_dimensions, value.dynamic_array, value.array_of_vectors);
    end

    function value = read(self, instream)
      fields = self.read_(instream);
      value = test_model.RecordWithArrays(default_array=fields{1}, default_array_with_empty_dimension=fields{2}, rank_1_array=fields{3}, rank_2_array=fields{4}, rank_2_array_with_named_dimensions=fields{5}, rank_2_fixed_array=fields{6}, rank_2_fixed_array_with_named_dimensions=fields{7}, dynamic_array=fields{8}, array_of_vectors=fields{9});
    end
  end
end
