% Binary reader for the StreamsOfUnions protocol
classdef StreamsOfUnionsReader < yardl.binary.BinaryProtocolReader & test_model.StreamsOfUnionsReaderBase
  methods
    function obj = StreamsOfUnionsReader(filename)
      obj@test_model.StreamsOfUnionsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.StreamsOfUnionsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_int_or_simple_record_(obj)
      r = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('test_model.Int32OrSimpleRecord', {yardl.binary.Int32Serializer, test_model.binary.SimpleRecordSerializer()}, {@test_model.Int32OrSimpleRecord.Int32, @test_model.Int32OrSimpleRecord.SimpleRecord}));
      value = r.read(obj.stream_);
    end

    function value = read_nullable_int_or_simple_record_(obj)
      r = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('test_model.Int32OrSimpleRecord', {yardl.binary.NoneSerializer, yardl.binary.Int32Serializer, test_model.binary.SimpleRecordSerializer()}, {yardl.None, @test_model.Int32OrSimpleRecord.Int32, @test_model.Int32OrSimpleRecord.SimpleRecord}));
      value = r.read(obj.stream_);
    end
  end
end
