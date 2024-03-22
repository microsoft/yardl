classdef RecordWithNDArraysSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithNDArraysSerializer()
      field_serializers{1} = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      field_serializers{2} = yardl.binary.NDArraySerializer(test_model.binary.SimpleRecordSerializer(), 2);
      field_serializers{3} = yardl.binary.NDArraySerializer(test_model.binary.RecordWithVlensSerializer(), 2);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithNDArrays', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithNDArrays'));
      obj.write_(outstream, value.ints, value.fixed_simple_record_array, value.fixed_record_with_vlens_array)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithNDArrays(field_values{:});
    end
  end
end
