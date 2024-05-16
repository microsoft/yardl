% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef StreamsOfUnionsReader < yardl.binary.BinaryProtocolReader & test_model.StreamsOfUnionsReaderBase
  % Binary reader for the StreamsOfUnions protocol
  properties (Access=protected)
    int_or_simple_record_serializer
    nullable_int_or_simple_record_serializer
  end

  methods
    function self = StreamsOfUnionsReader(filename)
      self@test_model.StreamsOfUnionsReaderBase();
      self@yardl.binary.BinaryProtocolReader(filename, test_model.StreamsOfUnionsReaderBase.schema);
      self.int_or_simple_record_serializer = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('test_model.Int32OrSimpleRecord', {yardl.binary.Int32Serializer, test_model.binary.SimpleRecordSerializer()}, {@test_model.Int32OrSimpleRecord.Int32, @test_model.Int32OrSimpleRecord.SimpleRecord}));
      self.nullable_int_or_simple_record_serializer = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('test_model.Int32OrSimpleRecord', {yardl.binary.NoneSerializer, yardl.binary.Int32Serializer, test_model.binary.SimpleRecordSerializer()}, {yardl.None, @test_model.Int32OrSimpleRecord.Int32, @test_model.Int32OrSimpleRecord.SimpleRecord}));
    end
  end

  methods (Access=protected)
    function more = has_int_or_simple_record_(self)
      more = self.int_or_simple_record_serializer.hasnext(self.stream_);
    end

    function value = read_int_or_simple_record_(self)
      value = self.int_or_simple_record_serializer.read(self.stream_);
    end

    function more = has_nullable_int_or_simple_record_(self)
      more = self.nullable_int_or_simple_record_serializer.hasnext(self.stream_);
    end

    function value = read_nullable_int_or_simple_record_(self)
      value = self.nullable_int_or_simple_record_serializer.read(self.stream_);
    end
  end
end
