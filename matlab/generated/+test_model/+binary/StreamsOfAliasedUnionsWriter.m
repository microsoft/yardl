% Binary writer for the StreamsOfAliasedUnions protocol
classdef StreamsOfAliasedUnionsWriter < yardl.binary.BinaryProtocolWriter & test_model.StreamsOfAliasedUnionsWriterBase
  methods
    function obj = StreamsOfAliasedUnionsWriter(filename)
      obj@test_model.StreamsOfAliasedUnionsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.StreamsOfAliasedUnionsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_int_or_simple_record_(obj, value)
      w = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('test_model.AliasedIntOrSimpleRecord', {yardl.binary.Int32Serializer, test_model.binary.SimpleRecordSerializer()}, {@test_model.AliasedIntOrSimpleRecord.Int32, @test_model.AliasedIntOrSimpleRecord.SimpleRecord}));
      w.write(obj.stream_, value);
    end

    function write_nullable_int_or_simple_record_(obj, value)
      w = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('test_model.AliasedNullableIntSimpleRecord', {yardl.binary.NoneSerializer, yardl.binary.Int32Serializer, test_model.binary.SimpleRecordSerializer()}, {yardl.None, @test_model.AliasedNullableIntSimpleRecord.Int32, @test_model.AliasedNullableIntSimpleRecord.SimpleRecord}));
      w.write(obj.stream_, value);
    end
  end
end
