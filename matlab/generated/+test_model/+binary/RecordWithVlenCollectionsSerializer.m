classdef RecordWithVlenCollectionsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithVlenCollectionsSerializer()
      field_serializers{1} = yardl.binary.VectorSerializer(yardl.binary.Int32Serializer);
      field_serializers{2} = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithVlenCollections', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithVlenCollections'));
      obj.write_(outstream, value.vector, value.array)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithVlenCollections(field_values{:});
    end
  end
end
