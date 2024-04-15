% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef BenchmarkFloat256x256Reader < yardl.binary.BinaryProtocolReader & test_model.BenchmarkFloat256x256ReaderBase
  % Binary reader for the BenchmarkFloat256x256 protocol
  properties (Access=protected)
    float256x256_serializer
  end

  methods
    function self = BenchmarkFloat256x256Reader(filename)
      self@test_model.BenchmarkFloat256x256ReaderBase();
      self@yardl.binary.BinaryProtocolReader(filename, test_model.BenchmarkFloat256x256ReaderBase.schema);
      self.float256x256_serializer = yardl.binary.StreamSerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Float32Serializer, [256, 256]));
    end
  end

  methods (Access=protected)
    function more = has_float256x256_(self)
      more = self.float256x256_serializer.hasnext(self.stream_);
    end

    function value = read_float256x256_(self)
      value = self.float256x256_serializer.read(self.stream_);
    end
  end
end
