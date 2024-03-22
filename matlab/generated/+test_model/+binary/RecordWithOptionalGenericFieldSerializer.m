classdef RecordWithOptionalGenericFieldSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithOptionalGenericFieldSerializer(t_serializer)
      field_serializers{1} = yardl.binary.OptionalSerializer(t_serializer);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithOptionalGenericField', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithOptionalGenericField'));
      obj.write_(outstream, value.v)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithOptionalGenericField(field_values{:});
    end
  end
end
