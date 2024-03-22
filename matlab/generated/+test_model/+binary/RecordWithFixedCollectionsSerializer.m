classdef RecordWithFixedCollectionsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithFixedCollectionsSerializer()
      field_serializers{1} = yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 3);
      field_serializers{2} = yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3, 2]);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithFixedCollections', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithFixedCollections'));
      obj.write_(outstream, value.fixed_vector, value.fixed_array)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithFixedCollections(field_values{:});
    end
  end
end
