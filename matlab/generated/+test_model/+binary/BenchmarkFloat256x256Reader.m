% Binary reader for the BenchmarkFloat256x256 protocol
classdef BenchmarkFloat256x256Reader < yardl.binary.BinaryProtocolReader & test_model.BenchmarkFloat256x256ReaderBase
  methods
    function obj = BenchmarkFloat256x256Reader(filename)
      obj@test_model.BenchmarkFloat256x256ReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.BenchmarkFloat256x256ReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_float256x256_(obj)
      r = yardl.binary.StreamSerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Float32Serializer, [256, 256]));
      value = r.read(obj.stream_);
    end
  end
end
