% Binary reader for the StreamsOfAliasedUnions protocol
classdef StreamsOfAliasedUnionsReader < yardl.binary.BinaryProtocolReader & test_model.StreamsOfAliasedUnionsReaderBase
  methods
    function obj = StreamsOfAliasedUnionsReader(filename)
      obj@test_model.StreamsOfAliasedUnionsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.StreamsOfAliasedUnionsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_int_or_simple_record_(obj)
      r = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('test_model.AliasedIntOrSimpleRecord', {yardl.binary.Int32Serializer, test_model.binary.SimpleRecordSerializer()}, {@test_model.AliasedIntOrSimpleRecord.Int32, @test_model.AliasedIntOrSimpleRecord.SimpleRecord}));
      value = r.read(obj.stream_);
    end

    function value = read_nullable_int_or_simple_record_(obj)
      r = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('test_model.AliasedNullableIntSimpleRecord', {yardl.binary.NoneSerializer, yardl.binary.Int32Serializer, test_model.binary.SimpleRecordSerializer()}, {yardl.None, @test_model.AliasedNullableIntSimpleRecord.Int32, @test_model.AliasedNullableIntSimpleRecord.SimpleRecord}));
      value = r.read(obj.stream_);
    end
  end
end
