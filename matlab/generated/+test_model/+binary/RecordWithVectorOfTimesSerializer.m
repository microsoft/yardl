classdef RecordWithVectorOfTimesSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithVectorOfTimesSerializer()
      field_serializers{1} = yardl.binary.VectorSerializer(yardl.binary.TimeSerializer);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithVectorOfTimes', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithVectorOfTimes'));
      obj.write_(outstream, value.times)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithVectorOfTimes(field_values{:});
    end
  end
end
