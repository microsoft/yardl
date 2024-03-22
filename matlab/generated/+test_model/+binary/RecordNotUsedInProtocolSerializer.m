classdef RecordNotUsedInProtocolSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordNotUsedInProtocolSerializer()
      field_serializers{1} = yardl.binary.UnionSerializer('test_model.GenericUnion3', {yardl.binary.Int32Serializer, yardl.binary.Float32Serializer, yardl.binary.StringSerializer}, {@test_model.GenericUnion3.T, @test_model.GenericUnion3.U, @test_model.GenericUnion3.V});
      field_serializers{2} = yardl.binary.UnionSerializer('test_model.GenericUnion3Alternate', {yardl.binary.Int32Serializer, yardl.binary.Float32Serializer, yardl.binary.StringSerializer}, {@test_model.GenericUnion3Alternate.U, @test_model.GenericUnion3Alternate.V, @test_model.GenericUnion3Alternate.W});
      obj@yardl.binary.RecordSerializer('test_model.RecordNotUsedInProtocol', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordNotUsedInProtocol'));
      obj.write_(outstream, value.u1, value.u2)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordNotUsedInProtocol(field_values{:});
    end
  end
end
