classdef RecordWithGenericMapsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithGenericMapsSerializer(t_serializer, u_serializer)
      field_serializers{1} = yardl.binary.MapSerializer(t_serializer, u_serializer);
      field_serializers{2} = yardl.binary.MapSerializer(t_serializer, u_serializer);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithGenericMaps', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithGenericMaps'));
      obj.write_(outstream, value.m, value.am)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithGenericMaps(field_values{:});
    end
  end
end
