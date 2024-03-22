% Binary reader for the Unions protocol
classdef UnionsReader < yardl.binary.BinaryProtocolReader & test_model.UnionsReaderBase
  methods
    function obj = UnionsReader(filename)
      obj@test_model.UnionsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.UnionsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_int_or_simple_record_(obj)
      r = yardl.binary.UnionSerializer('test_model.Int32OrSimpleRecord', {yardl.binary.Int32Serializer, test_model.binary.SimpleRecordSerializer()}, {@test_model.Int32OrSimpleRecord.Int32, @test_model.Int32OrSimpleRecord.SimpleRecord});
      value = r.read(obj.stream_);
    end

    function value = read_int_or_record_with_vlens_(obj)
      r = yardl.binary.UnionSerializer('test_model.Int32OrRecordWithVlens', {yardl.binary.Int32Serializer, test_model.binary.RecordWithVlensSerializer()}, {@test_model.Int32OrRecordWithVlens.Int32, @test_model.Int32OrRecordWithVlens.RecordWithVlens});
      value = r.read(obj.stream_);
    end

    function value = read_monosotate_or_int_or_simple_record_(obj)
      r = yardl.binary.UnionSerializer('test_model.Int32OrSimpleRecord', {yardl.binary.NoneSerializer, yardl.binary.Int32Serializer, test_model.binary.SimpleRecordSerializer()}, {yardl.None, @test_model.Int32OrSimpleRecord.Int32, @test_model.Int32OrSimpleRecord.SimpleRecord});
      value = r.read(obj.stream_);
    end

    function value = read_record_with_unions_(obj)
      r = basic_types.binary.RecordWithUnionsSerializer();
      value = r.read(obj.stream_);
    end
  end
end
