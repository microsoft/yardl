classdef RecordWithFixedArraysSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithFixedArraysSerializer()
      field_serializers{1} = yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3, 2]);
      field_serializers{2} = yardl.binary.FixedNDArraySerializer(test_model.binary.SimpleRecordSerializer(), [2, 3]);
      field_serializers{3} = yardl.binary.FixedNDArraySerializer(test_model.binary.RecordWithVlensSerializer(), [2, 2]);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithFixedArrays', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithFixedArrays'));
      obj.write_(outstream, value.ints, value.fixed_simple_record_array, value.fixed_record_with_vlens_array)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithFixedArrays(field_values{:});
    end
  end
end