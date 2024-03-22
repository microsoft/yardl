% Binary writer for the StreamsOfUnions protocol
classdef StreamsOfUnionsWriter < yardl.binary.BinaryProtocolWriter & test_model.StreamsOfUnionsWriterBase
  methods
    function obj = StreamsOfUnionsWriter(filename)
      obj@test_model.StreamsOfUnionsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.StreamsOfUnionsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_int_or_simple_record_(obj, value)
      w = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('test_model.Int32OrSimpleRecord', {yardl.binary.Int32Serializer, test_model.binary.SimpleRecordSerializer()}, {@test_model.Int32OrSimpleRecord.Int32, @test_model.Int32OrSimpleRecord.SimpleRecord}));
      w.write(obj.stream_, value);
    end

    function write_nullable_int_or_simple_record_(obj, value)
      w = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('test_model.Int32OrSimpleRecord', {yardl.binary.NoneSerializer, yardl.binary.Int32Serializer, test_model.binary.SimpleRecordSerializer()}, {yardl.None, @test_model.Int32OrSimpleRecord.Int32, @test_model.Int32OrSimpleRecord.SimpleRecord}));
      w.write(obj.stream_, value);
    end
  end
end