% Binary reader for the BenchmarkSmallRecord protocol
classdef BenchmarkSmallRecordReader < yardl.binary.BinaryProtocolReader & test_model.BenchmarkSmallRecordReaderBase
  methods
    function obj = BenchmarkSmallRecordReader(filename)
      obj@test_model.BenchmarkSmallRecordReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.BenchmarkSmallRecordReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_small_record_(obj)
      r = yardl.binary.StreamSerializer(test_model.binary.SmallBenchmarkRecordSerializer());
      value = r.read(obj.stream_);
    end
  end
end