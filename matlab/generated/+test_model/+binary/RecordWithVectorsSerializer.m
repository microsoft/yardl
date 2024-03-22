classdef RecordWithVectorsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithVectorsSerializer()
      field_serializers{1} = yardl.binary.VectorSerializer(yardl.binary.Int32Serializer);
      field_serializers{2} = yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 3);
      field_serializers{3} = yardl.binary.VectorSerializer(yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 2));
      obj@yardl.binary.RecordSerializer('test_model.RecordWithVectors', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithVectors'));
      obj.write_(outstream, value.default_vector, value.default_vector_fixed_length, value.vector_of_vectors)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithVectors(field_values{:});
    end
  end
end