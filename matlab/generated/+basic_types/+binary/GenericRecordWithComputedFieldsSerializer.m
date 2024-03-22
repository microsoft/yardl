classdef GenericRecordWithComputedFieldsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = GenericRecordWithComputedFieldsSerializer(t0_serializer, t1_serializer)
      field_serializers{1} = yardl.binary.UnionSerializer('basic_types.T0OrT1', {t0_serializer, t1_serializer}, {@basic_types.T0OrT1.T0, @basic_types.T0OrT1.T1});
      obj@yardl.binary.RecordSerializer('basic_types.GenericRecordWithComputedFields', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'basic_types.GenericRecordWithComputedFields'));
      obj.write_(outstream, value.f1)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = basic_types.GenericRecordWithComputedFields(field_values{:});
    end
  end
end
