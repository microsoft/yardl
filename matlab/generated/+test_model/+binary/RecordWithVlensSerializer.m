classdef RecordWithVlensSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithVlensSerializer()
      field_serializers{1} = yardl.binary.VectorSerializer(test_model.binary.SimpleRecordSerializer());
      field_serializers{2} = yardl.binary.Int32Serializer;
      field_serializers{3} = yardl.binary.Int32Serializer;
      obj@yardl.binary.RecordSerializer('test_model.RecordWithVlens', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithVlens'));
      obj.write_(outstream, value.a, value.b, value.c)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithVlens(field_values{:});
    end
  end
end
