% Binary writer for the BenchmarkFloat256x256 protocol
classdef BenchmarkFloat256x256Writer < yardl.binary.BinaryProtocolWriter & test_model.BenchmarkFloat256x256WriterBase
  methods
    function obj = BenchmarkFloat256x256Writer(filename)
      obj@test_model.BenchmarkFloat256x256WriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.BenchmarkFloat256x256WriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_float256x256_(obj, value)
      w = yardl.binary.StreamSerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Float32Serializer, [256, 256]));
      w.write(obj.stream_, value);
    end
  end
end
