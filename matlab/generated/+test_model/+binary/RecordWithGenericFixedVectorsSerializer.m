% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithGenericFixedVectorsSerializer < yardl.binary.RecordSerializer
  methods
    function self = RecordWithGenericFixedVectorsSerializer(t_serializer)
      field_serializers{1} = yardl.binary.FixedVectorSerializer(t_serializer, 3);
      field_serializers{2} = yardl.binary.FixedVectorSerializer(t_serializer, 3);
      self@yardl.binary.RecordSerializer('test_model.RecordWithGenericFixedVectors', field_serializers);
    end

    function write(self, outstream, value)
      arguments
        self
        outstream (1,1) yardl.binary.CodedOutputStream
        value (1,1) test_model.RecordWithGenericFixedVectors
      end
      self.write_(outstream, value.fv, value.afv);
    end

    function value = read(self, instream)
      fields = self.read_(instream);
      value = test_model.RecordWithGenericFixedVectors(fv=fields{1}, afv=fields{2});
    end
  end
end
