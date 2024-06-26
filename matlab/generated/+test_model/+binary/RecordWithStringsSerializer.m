% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithStringsSerializer < yardl.binary.RecordSerializer
  methods
    function self = RecordWithStringsSerializer()
      field_serializers{1} = yardl.binary.StringSerializer;
      field_serializers{2} = yardl.binary.StringSerializer;
      self@yardl.binary.RecordSerializer('test_model.RecordWithStrings', field_serializers);
    end

    function write(self, outstream, value)
      arguments
        self
        outstream (1,1) yardl.binary.CodedOutputStream
        value (1,1) test_model.RecordWithStrings
      end
      self.write_(outstream, value.a, value.b);
    end

    function value = read(self, instream)
      fields = self.read_(instream);
      value = test_model.RecordWithStrings(a=fields{1}, b=fields{2});
    end
  end
end
