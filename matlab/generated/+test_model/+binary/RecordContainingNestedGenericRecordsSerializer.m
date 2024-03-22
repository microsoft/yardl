classdef RecordContainingNestedGenericRecordsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordContainingNestedGenericRecordsSerializer()
      field_serializers{1} = test_model.binary.RecordWithOptionalGenericFieldSerializer(yardl.binary.StringSerializer);
      field_serializers{2} = test_model.binary.RecordWithAliasedOptionalGenericFieldSerializer(yardl.binary.StringSerializer);
      field_serializers{3} = test_model.binary.RecordWithOptionalGenericUnionFieldSerializer(yardl.binary.StringSerializer, yardl.binary.Int32Serializer);
      field_serializers{4} = test_model.binary.RecordWithAliasedOptionalGenericUnionFieldSerializer(yardl.binary.StringSerializer, yardl.binary.Int32Serializer);
      field_serializers{5} = test_model.binary.RecordContainingGenericRecordsSerializer(yardl.binary.StringSerializer, yardl.binary.Int32Serializer);
      obj@yardl.binary.RecordSerializer('test_model.RecordContainingNestedGenericRecords', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordContainingNestedGenericRecords'));
      obj.write_(outstream, value.f1, value.f1a, value.f2, value.f2a, value.nested)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordContainingNestedGenericRecords(field_values{:});
    end
  end
end