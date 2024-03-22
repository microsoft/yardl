classdef RecordWithKeywordFieldsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithKeywordFieldsSerializer()
      field_serializers{1} = yardl.binary.StringSerializer;
      field_serializers{2} = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      field_serializers{3} = yardl.binary.EnumSerializer('test_model.EnumWithKeywordSymbols', @test_model.EnumWithKeywordSymbols, yardl.binary.Int32Serializer);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithKeywordFields', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithKeywordFields'));
      obj.write_(outstream, value.int, value.sizeof, value.if)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithKeywordFields(field_values{:});
    end
  end
end
