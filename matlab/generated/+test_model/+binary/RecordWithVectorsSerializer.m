% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithVectorsSerializer < yardl.binary.RecordSerializer
  methods
    function self = RecordWithVectorsSerializer()
      field_serializers{1} = yardl.binary.VectorSerializer(yardl.binary.Int32Serializer);
      field_serializers{2} = yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 3);
      field_serializers{3} = yardl.binary.VectorSerializer(yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 2));
      self@yardl.binary.RecordSerializer('test_model.RecordWithVectors', field_serializers);
    end

    function write(self, outstream, value)
      arguments
        self
        outstream (1,1) yardl.binary.CodedOutputStream
        value (1,1) test_model.RecordWithVectors
      end
      self.write_(outstream, value.default_vector, value.default_vector_fixed_length, value.vector_of_vectors)
    end

    function value = read(self, instream)
      field_values = self.read_(instream);
      value = test_model.RecordWithVectors(field_values{:});
    end
  end
end
