% Binary reader for the BenchmarkInt256x256 protocol
classdef BenchmarkInt256x256Reader < yardl.binary.BinaryProtocolReader & test_model.BenchmarkInt256x256ReaderBase
  methods
    function obj = BenchmarkInt256x256Reader(filename)
      obj@test_model.BenchmarkInt256x256ReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.BenchmarkInt256x256ReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_int256x256_(obj)
      r = yardl.binary.StreamSerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [256, 256]));
      value = r.read(obj.stream_);
    end
  end
end
