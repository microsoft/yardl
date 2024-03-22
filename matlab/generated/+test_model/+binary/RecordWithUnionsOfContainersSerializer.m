classdef RecordWithUnionsOfContainersSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithUnionsOfContainersSerializer()
      field_serializers{1} = yardl.binary.UnionSerializer('test_model.MapOrScalar', {yardl.binary.MapSerializer(yardl.binary.StringSerializer, yardl.binary.Int32Serializer), yardl.binary.Int32Serializer}, {@test_model.MapOrScalar.Map, @test_model.MapOrScalar.Scalar});
      field_serializers{2} = yardl.binary.UnionSerializer('test_model.VectorOrScalar', {yardl.binary.VectorSerializer(yardl.binary.Int32Serializer), yardl.binary.Int32Serializer}, {@test_model.VectorOrScalar.Vector, @test_model.VectorOrScalar.Scalar});
      field_serializers{3} = yardl.binary.UnionSerializer('test_model.ArrayOrScalar', {yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer), yardl.binary.Int32Serializer}, {@test_model.ArrayOrScalar.Array, @test_model.ArrayOrScalar.Scalar});
      obj@yardl.binary.RecordSerializer('test_model.RecordWithUnionsOfContainers', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithUnionsOfContainers'));
      obj.write_(outstream, value.map_or_scalar, value.vector_or_scalar, value.array_or_scalar)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithUnionsOfContainers(field_values{:});
    end
  end
end