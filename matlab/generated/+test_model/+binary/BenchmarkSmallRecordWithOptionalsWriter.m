% Binary writer for the BenchmarkSmallRecordWithOptionals protocol
classdef BenchmarkSmallRecordWithOptionalsWriter < yardl.binary.BinaryProtocolWriter & test_model.BenchmarkSmallRecordWithOptionalsWriterBase
  methods
    function obj = BenchmarkSmallRecordWithOptionalsWriter(filename)
      obj@test_model.BenchmarkSmallRecordWithOptionalsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.BenchmarkSmallRecordWithOptionalsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_small_record_(obj, value)
      w = yardl.binary.StreamSerializer(test_model.binary.SimpleEncodingCountersSerializer());
      w.write(obj.stream_, value);
    end
  end
end
