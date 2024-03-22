classdef RecordWithGenericVectorsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithGenericVectorsSerializer(t_serializer)
      field_serializers{1} = yardl.binary.VectorSerializer(t_serializer);
      field_serializers{2} = yardl.binary.VectorSerializer(t_serializer);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithGenericVectors', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithGenericVectors'));
      obj.write_(outstream, value.v, value.av)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithGenericVectors(field_values{:});
    end
  end
end
