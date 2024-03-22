classdef RecordWithOptionalGenericUnionFieldSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithOptionalGenericUnionFieldSerializer(u_serializer, v_serializer)
      field_serializers{1} = yardl.binary.UnionSerializer('test_model.UOrV', {yardl.binary.NoneSerializer, u_serializer, v_serializer}, {yardl.None, @test_model.UOrV.U, @test_model.UOrV.V});
      obj@yardl.binary.RecordSerializer('test_model.RecordWithOptionalGenericUnionField', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithOptionalGenericUnionField'));
      obj.write_(outstream, value.v)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithOptionalGenericUnionField(field_values{:});
    end
  end
end
