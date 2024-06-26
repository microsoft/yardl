% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef BenchmarkFloat256x256Writer < yardl.binary.BinaryProtocolWriter & test_model.BenchmarkFloat256x256WriterBase
  % Binary writer for the BenchmarkFloat256x256 protocol
  properties (Access=protected)
    float256x256_serializer
  end

  methods
    function self = BenchmarkFloat256x256Writer(filename)
      self@test_model.BenchmarkFloat256x256WriterBase();
      self@yardl.binary.BinaryProtocolWriter(filename, test_model.BenchmarkFloat256x256WriterBase.schema);
      self.float256x256_serializer = yardl.binary.StreamSerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Float32Serializer, [256, 256]));
    end
  end

  methods (Access=protected)
    function write_float256x256_(self, value)
      self.float256x256_serializer.write(self.stream_, value);
    end
  end
end
