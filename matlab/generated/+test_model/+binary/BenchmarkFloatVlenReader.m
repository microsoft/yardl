% Binary reader for the BenchmarkFloatVlen protocol
classdef BenchmarkFloatVlenReader < yardl.binary.BinaryProtocolReader & test_model.BenchmarkFloatVlenReaderBase
  methods
    function obj = BenchmarkFloatVlenReader(filename)
      obj@test_model.BenchmarkFloatVlenReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.BenchmarkFloatVlenReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_float_array_(obj)
      r = yardl.binary.StreamSerializer(yardl.binary.NDArraySerializer(yardl.binary.Float32Serializer, 2));
      value = r.read(obj.stream_);
    end
  end
end
