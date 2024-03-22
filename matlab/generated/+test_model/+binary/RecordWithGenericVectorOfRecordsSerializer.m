classdef RecordWithGenericVectorOfRecordsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithGenericVectorOfRecordsSerializer(t_serializer, u_serializer)
      field_serializers{1} = yardl.binary.VectorSerializer(yardl.binary.VectorSerializer(test_model.binary.GenericRecordSerializer(t_serializer, u_serializer)));
      obj@yardl.binary.RecordSerializer('test_model.RecordWithGenericVectorOfRecords', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithGenericVectorOfRecords'));
      obj.write_(outstream, value.v)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithGenericVectorOfRecords(field_values{:});
    end
  end
end
