% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithVectorOfTimesSerializer < yardl.binary.RecordSerializer
  methods
    function self = RecordWithVectorOfTimesSerializer()
      field_serializers{1} = yardl.binary.VectorSerializer(yardl.binary.TimeSerializer);
      self@yardl.binary.RecordSerializer('test_model.RecordWithVectorOfTimes', field_serializers);
    end

    function write(self, outstream, value)
      arguments
        self
        outstream (1,1) yardl.binary.CodedOutputStream
        value (1,1) test_model.RecordWithVectorOfTimes
      end
      self.write_(outstream, value.times)
    end

    function value = read(self, instream)
      field_values = self.read_(instream);
      value = test_model.RecordWithVectorOfTimes(field_values{:});
    end
  end
end
