classdef RecordWithGenericFixedVectorsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithGenericFixedVectorsSerializer(t_serializer)
      field_serializers{1} = yardl.binary.FixedVectorSerializer(t_serializer, 3);
      field_serializers{2} = yardl.binary.FixedVectorSerializer(t_serializer, 3);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithGenericFixedVectors', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithGenericFixedVectors'));
      obj.write_(outstream, value.fv, value.afv)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithGenericFixedVectors(field_values{:});
    end
  end
end
