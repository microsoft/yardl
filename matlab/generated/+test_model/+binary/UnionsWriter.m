% Binary writer for the Unions protocol
classdef UnionsWriter < yardl.binary.BinaryProtocolWriter & test_model.UnionsWriterBase
  methods
    function obj = UnionsWriter(filename)
      obj@test_model.UnionsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.UnionsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_int_or_simple_record_(obj, value)
      w = yardl.binary.UnionSerializer('test_model.Int32OrSimpleRecord', {yardl.binary.Int32Serializer, test_model.binary.SimpleRecordSerializer()}, {@test_model.Int32OrSimpleRecord.Int32, @test_model.Int32OrSimpleRecord.SimpleRecord});
      w.write(obj.stream_, value);
    end

    function write_int_or_record_with_vlens_(obj, value)
      w = yardl.binary.UnionSerializer('test_model.Int32OrRecordWithVlens', {yardl.binary.Int32Serializer, test_model.binary.RecordWithVlensSerializer()}, {@test_model.Int32OrRecordWithVlens.Int32, @test_model.Int32OrRecordWithVlens.RecordWithVlens});
      w.write(obj.stream_, value);
    end

    function write_monosotate_or_int_or_simple_record_(obj, value)
      w = yardl.binary.UnionSerializer('test_model.Int32OrSimpleRecord', {yardl.binary.NoneSerializer, yardl.binary.Int32Serializer, test_model.binary.SimpleRecordSerializer()}, {yardl.None, @test_model.Int32OrSimpleRecord.Int32, @test_model.Int32OrSimpleRecord.SimpleRecord});
      w.write(obj.stream_, value);
    end

    function write_record_with_unions_(obj, value)
      w = basic_types.binary.RecordWithUnionsSerializer();
      w.write(obj.stream_, value);
    end
  end
end
