classdef RecordWithOptionalFieldsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithOptionalFieldsSerializer()
      field_serializers{1} = yardl.binary.OptionalSerializer(yardl.binary.Int32Serializer);
      field_serializers{2} = yardl.binary.OptionalSerializer(yardl.binary.Int32Serializer);
      field_serializers{3} = yardl.binary.OptionalSerializer(yardl.binary.TimeSerializer);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithOptionalFields', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithOptionalFields'));
      obj.write_(outstream, value.optional_int, value.optional_int_alternate_syntax, value.optional_time)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithOptionalFields(field_values{:});
    end
  end
end