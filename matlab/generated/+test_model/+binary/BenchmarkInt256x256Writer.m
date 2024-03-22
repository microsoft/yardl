% Binary writer for the BenchmarkInt256x256 protocol
classdef BenchmarkInt256x256Writer < yardl.binary.BinaryProtocolWriter & test_model.BenchmarkInt256x256WriterBase
  methods
    function obj = BenchmarkInt256x256Writer(filename)
      obj@test_model.BenchmarkInt256x256WriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.BenchmarkInt256x256WriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_int256x256_(obj, value)
      w = yardl.binary.StreamSerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [256, 256]));
      w.write(obj.stream_, value);
    end
  end
end