% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithAliasedOptionalGenericFieldSerializer < yardl.binary.RecordSerializer
  methods
    function self = RecordWithAliasedOptionalGenericFieldSerializer(t_serializer)
      field_serializers{1} = yardl.binary.OptionalSerializer(t_serializer);
      self@yardl.binary.RecordSerializer('test_model.RecordWithAliasedOptionalGenericField', field_serializers);
    end

    function write(self, outstream, value)
      arguments
        self
        outstream (1,1) yardl.binary.CodedOutputStream
        value (1,1) test_model.RecordWithAliasedOptionalGenericField
      end
      self.write_(outstream, value.v);
    end

    function value = read(self, instream)
      fields = self.read_(instream);
      value = test_model.RecordWithAliasedOptionalGenericField(v=fields{1});
    end
  end
end
