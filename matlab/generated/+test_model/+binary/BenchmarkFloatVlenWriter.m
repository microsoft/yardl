% Binary writer for the BenchmarkFloatVlen protocol
classdef BenchmarkFloatVlenWriter < yardl.binary.BinaryProtocolWriter & test_model.BenchmarkFloatVlenWriterBase
  methods
    function obj = BenchmarkFloatVlenWriter(filename)
      obj@test_model.BenchmarkFloatVlenWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.BenchmarkFloatVlenWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_float_array_(obj, value)
      w = yardl.binary.StreamSerializer(yardl.binary.NDArraySerializer(yardl.binary.Float32Serializer, 2));
      w.write(obj.stream_, value);
    end
  end
end
