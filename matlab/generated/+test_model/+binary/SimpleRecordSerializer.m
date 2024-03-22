classdef SimpleRecordSerializer < yardl.binary.RecordSerializer
  methods
    function obj = SimpleRecordSerializer()
      field_serializers{1} = yardl.binary.Int32Serializer;
      field_serializers{2} = yardl.binary.Int32Serializer;
      field_serializers{3} = yardl.binary.Int32Serializer;
      obj@yardl.binary.RecordSerializer('test_model.SimpleRecord', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.SimpleRecord'));
      obj.write_(outstream, value.x, value.y, value.z)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.SimpleRecord(field_values{:});
    end
  end
end
