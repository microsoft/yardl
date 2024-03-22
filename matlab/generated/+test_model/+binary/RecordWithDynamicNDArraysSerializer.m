classdef RecordWithDynamicNDArraysSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithDynamicNDArraysSerializer()
      field_serializers{1} = yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer);
      field_serializers{2} = yardl.binary.DynamicNDArraySerializer(test_model.binary.SimpleRecordSerializer());
      field_serializers{3} = yardl.binary.DynamicNDArraySerializer(test_model.binary.RecordWithVlensSerializer());
      obj@yardl.binary.RecordSerializer('test_model.RecordWithDynamicNDArrays', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithDynamicNDArrays'));
      obj.write_(outstream, value.ints, value.simple_record_array, value.record_with_vlens_array)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithDynamicNDArrays(field_values{:});
    end
  end
end