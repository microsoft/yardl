% Binary writer for the BenchmarkSmallRecord protocol
classdef BenchmarkSmallRecordWriter < yardl.binary.BinaryProtocolWriter & test_model.BenchmarkSmallRecordWriterBase
  methods
    function obj = BenchmarkSmallRecordWriter(filename)
      obj@test_model.BenchmarkSmallRecordWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.BenchmarkSmallRecordWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_small_record_(obj, value)
      w = yardl.binary.StreamSerializer(test_model.binary.SmallBenchmarkRecordSerializer());
      w.write(obj.stream_, value);
    end
  end
end
