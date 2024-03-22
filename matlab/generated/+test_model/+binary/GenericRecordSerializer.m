classdef GenericRecordSerializer < yardl.binary.RecordSerializer
  methods
    function obj = GenericRecordSerializer(t1_serializer, t2_serializer)
      field_serializers{1} = t1_serializer;
      field_serializers{2} = t2_serializer;
      field_serializers{3} = yardl.binary.VectorSerializer(t1_serializer);
      field_serializers{4} = yardl.binary.NDArraySerializer(t2_serializer, 2);
      obj@yardl.binary.RecordSerializer('test_model.GenericRecord', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.GenericRecord'));
      obj.write_(outstream, value.scalar_1, value.scalar_2, value.vector_1, value.image_2)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.GenericRecord(field_values{:});
    end
  end
end
