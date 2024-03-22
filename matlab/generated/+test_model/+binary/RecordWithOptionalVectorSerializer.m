classdef RecordWithOptionalVectorSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithOptionalVectorSerializer()
      field_serializers{1} = yardl.binary.OptionalSerializer(yardl.binary.VectorSerializer(yardl.binary.Int32Serializer));
      obj@yardl.binary.RecordSerializer('test_model.RecordWithOptionalVector', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithOptionalVector'));
      obj.write_(outstream, value.optional_vector)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithOptionalVector(field_values{:});
    end
  end
end