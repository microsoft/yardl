classdef RecordWithFixedVectorsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithFixedVectorsSerializer()
      field_serializers{1} = yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 5);
      field_serializers{2} = yardl.binary.FixedVectorSerializer(test_model.binary.SimpleRecordSerializer(), 3);
      field_serializers{3} = yardl.binary.FixedVectorSerializer(test_model.binary.RecordWithVlensSerializer(), 2);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithFixedVectors', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithFixedVectors'));
      obj.write_(outstream, value.fixed_int_vector, value.fixed_simple_record_vector, value.fixed_record_with_vlens_vector)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithFixedVectors(field_values{:});
    end
  end
end
