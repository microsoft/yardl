% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef TupleSerializer < yardl.binary.RecordSerializer
  methods
    function self = TupleSerializer(t1_serializer, t2_serializer)
      field_serializers{1} = t1_serializer;
      field_serializers{2} = t2_serializer;
      self@yardl.binary.RecordSerializer('tuples.Tuple', field_serializers);
    end

    function write(self, outstream, value)
      arguments
        self
        outstream (1,1) yardl.binary.CodedOutputStream
        value (1,1) tuples.Tuple
      end
      self.write_(outstream, value.v1, value.v2);
    end

    function value = read(self, instream)
      fields = self.read_(instream);
      value = tuples.Tuple(v1=fields{1}, v2=fields{2});
    end
  end
end
