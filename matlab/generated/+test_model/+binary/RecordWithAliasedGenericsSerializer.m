% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithAliasedGenericsSerializer < yardl.binary.RecordSerializer
  methods
    function self = RecordWithAliasedGenericsSerializer()
      field_serializers{1} = tuples.binary.TupleSerializer(yardl.binary.StringSerializer, yardl.binary.StringSerializer);
      field_serializers{2} = tuples.binary.TupleSerializer(yardl.binary.StringSerializer, yardl.binary.StringSerializer);
      self@yardl.binary.RecordSerializer('test_model.RecordWithAliasedGenerics', field_serializers);
    end

    function write(self, outstream, value)
      arguments
        self
        outstream (1,1) yardl.binary.CodedOutputStream
        value (1,1) test_model.RecordWithAliasedGenerics
      end
      self.write_(outstream, value.my_strings, value.aliased_strings)
    end

    function value = read(self, instream)
      field_values = self.read_(instream);
      value = test_model.RecordWithAliasedGenerics(field_values{:});
    end
  end
end
