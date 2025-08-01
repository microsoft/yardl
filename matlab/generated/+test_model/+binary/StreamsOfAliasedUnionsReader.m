% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef StreamsOfAliasedUnionsReader < yardl.binary.BinaryProtocolReader & test_model.StreamsOfAliasedUnionsReaderBase
  % Binary reader for the StreamsOfAliasedUnions protocol
  properties (Access=protected)
    int_or_simple_record_serializer
    nullable_int_or_simple_record_serializer
  end

  methods
    function self = StreamsOfAliasedUnionsReader(filename, options)
      arguments
        filename (1,1) string
        options.skip_completed_check (1,1) logical = false
      end
      self@test_model.StreamsOfAliasedUnionsReaderBase(skip_completed_check=options.skip_completed_check);
      self@yardl.binary.BinaryProtocolReader(filename, test_model.StreamsOfAliasedUnionsReaderBase.schema);
      self.int_or_simple_record_serializer = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('test_model.AliasedIntOrSimpleRecord', {yardl.binary.Int32Serializer, test_model.binary.SimpleRecordSerializer()}, {@test_model.AliasedIntOrSimpleRecord.Int32, @test_model.AliasedIntOrSimpleRecord.SimpleRecord}));
      self.nullable_int_or_simple_record_serializer = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('test_model.AliasedNullableIntSimpleRecord', {yardl.binary.NoneSerializer, yardl.binary.Int32Serializer, test_model.binary.SimpleRecordSerializer()}, {yardl.None, @test_model.AliasedNullableIntSimpleRecord.Int32, @test_model.AliasedNullableIntSimpleRecord.SimpleRecord}));
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
