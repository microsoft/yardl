classdef RecordContainingGenericRecordsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordContainingGenericRecordsSerializer(a_serializer, b_serializer)
      field_serializers{1} = test_model.binary.RecordWithOptionalGenericFieldSerializer(a_serializer);
      field_serializers{2} = test_model.binary.RecordWithAliasedOptionalGenericFieldSerializer(a_serializer);
      field_serializers{3} = test_model.binary.RecordWithOptionalGenericUnionFieldSerializer(a_serializer, b_serializer);
      field_serializers{4} = test_model.binary.RecordWithAliasedOptionalGenericUnionFieldSerializer(a_serializer, b_serializer);
      field_serializers{5} = tuples.binary.TupleSerializer(a_serializer, b_serializer);
      field_serializers{6} = tuples.binary.TupleSerializer(a_serializer, b_serializer);
      field_serializers{7} = test_model.binary.RecordWithGenericVectorsSerializer(b_serializer);
      field_serializers{8} = test_model.binary.RecordWithGenericFixedVectorsSerializer(b_serializer);
      field_serializers{9} = test_model.binary.RecordWithGenericArraysSerializer(b_serializer);
      field_serializers{10} = test_model.binary.RecordWithGenericMapsSerializer(a_serializer, b_serializer);
      obj@yardl.binary.RecordSerializer('test_model.RecordContainingGenericRecords', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordContainingGenericRecords'));
      obj.write_(outstream, value.g1, value.g1a, value.g2, value.g2a, value.g3, value.g3a, value.g4, value.g5, value.g6, value.g7)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordContainingGenericRecords(field_values{:});
    end
  end
end