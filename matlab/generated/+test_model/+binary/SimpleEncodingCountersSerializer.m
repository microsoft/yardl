classdef SimpleEncodingCountersSerializer < yardl.binary.RecordSerializer
  methods
    function obj = SimpleEncodingCountersSerializer()
      field_serializers{1} = yardl.binary.OptionalSerializer(yardl.binary.Uint32Serializer);
      field_serializers{2} = yardl.binary.OptionalSerializer(yardl.binary.Uint32Serializer);
      field_serializers{3} = yardl.binary.OptionalSerializer(yardl.binary.Uint32Serializer);
      field_serializers{4} = yardl.binary.OptionalSerializer(yardl.binary.Uint32Serializer);
      obj@yardl.binary.RecordSerializer('test_model.SimpleEncodingCounters', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.SimpleEncodingCounters'));
      obj.write_(outstream, value.e1, value.e2, value.slice, value.repetition)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.SimpleEncodingCounters(field_values{:});
    end
  end
end
