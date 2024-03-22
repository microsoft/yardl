% Binary reader for the BenchmarkSmallRecordWithOptionals protocol
classdef BenchmarkSmallRecordWithOptionalsReader < yardl.binary.BinaryProtocolReader & test_model.BenchmarkSmallRecordWithOptionalsReaderBase
  methods
    function obj = BenchmarkSmallRecordWithOptionalsReader(filename)
      obj@test_model.BenchmarkSmallRecordWithOptionalsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.BenchmarkSmallRecordWithOptionalsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_small_record_(obj)
      r = yardl.binary.StreamSerializer(test_model.binary.SimpleEncodingCountersSerializer());
      value = r.read(obj.stream_);
    end
  end
end
