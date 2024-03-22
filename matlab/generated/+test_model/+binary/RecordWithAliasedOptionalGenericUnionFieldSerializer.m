classdef RecordWithAliasedOptionalGenericUnionFieldSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithAliasedOptionalGenericUnionFieldSerializer(u_serializer, v_serializer)
      field_serializers{1} = yardl.binary.UnionSerializer('test_model.AliasedMultiGenericOptional', {yardl.binary.NoneSerializer, u_serializer, v_serializer}, {yardl.None, @test_model.AliasedMultiGenericOptional.T, @test_model.AliasedMultiGenericOptional.U});
      obj@yardl.binary.RecordSerializer('test_model.RecordWithAliasedOptionalGenericUnionField', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithAliasedOptionalGenericUnionField'));
      obj.write_(outstream, value.v)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithAliasedOptionalGenericUnionField(field_values{:});
    end
  end
end
