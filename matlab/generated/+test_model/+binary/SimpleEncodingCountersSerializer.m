% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef SimpleEncodingCountersSerializer < yardl.binary.RecordSerializer
  methods
    function self = SimpleEncodingCountersSerializer()
      field_serializers{1} = yardl.binary.OptionalSerializer(yardl.binary.Uint32Serializer);
      field_serializers{2} = yardl.binary.OptionalSerializer(yardl.binary.Uint32Serializer);
      field_serializers{3} = yardl.binary.OptionalSerializer(yardl.binary.Uint32Serializer);
      field_serializers{4} = yardl.binary.OptionalSerializer(yardl.binary.Uint32Serializer);
      self@yardl.binary.RecordSerializer('test_model.SimpleEncodingCounters', field_serializers);
    end

    function write(self, outstream, value)
      arguments
        self
        outstream (1,1) yardl.binary.CodedOutputStream
        value (1,1) test_model.SimpleEncodingCounters
      end
      self.write_(outstream, value.e1, value.e2, value.slice, value.repetition);
    end

    function value = read(self, instream)
      fields = self.read_(instream);
      value = test_model.SimpleEncodingCounters(e1=fields{1}, e2=fields{2}, slice=fields{3}, repetition=fields{4});
    end
  end
end
