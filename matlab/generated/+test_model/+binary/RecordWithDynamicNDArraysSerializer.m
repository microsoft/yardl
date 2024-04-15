% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithDynamicNDArraysSerializer < yardl.binary.RecordSerializer
  methods
    function self = RecordWithDynamicNDArraysSerializer()
      field_serializers{1} = yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer);
      field_serializers{2} = yardl.binary.DynamicNDArraySerializer(test_model.binary.SimpleRecordSerializer());
      field_serializers{3} = yardl.binary.DynamicNDArraySerializer(test_model.binary.RecordWithVlensSerializer());
      self@yardl.binary.RecordSerializer('test_model.RecordWithDynamicNDArrays', field_serializers);
    end

    function write(self, outstream, value)
      arguments
        self
        outstream (1,1) yardl.binary.CodedOutputStream
        value (1,1) test_model.RecordWithDynamicNDArrays
      end
      self.write_(outstream, value.ints, value.simple_record_array, value.record_with_vlens_array)
    end

    function value = read(self, instream)
      field_values = self.read_(instream);
      value = test_model.RecordWithDynamicNDArrays(field_values{:});
    end
  end
end
