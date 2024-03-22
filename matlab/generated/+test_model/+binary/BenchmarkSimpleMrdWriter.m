% Binary writer for the BenchmarkSimpleMrd protocol
classdef BenchmarkSimpleMrdWriter < yardl.binary.BinaryProtocolWriter & test_model.BenchmarkSimpleMrdWriterBase
  methods
    function obj = BenchmarkSimpleMrdWriter(filename)
      obj@test_model.BenchmarkSimpleMrdWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.BenchmarkSimpleMrdWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_data_(obj, value)
      w = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('test_model.AcquisitionOrImage', {test_model.binary.SimpleAcquisitionSerializer(), yardl.binary.NDArraySerializer(yardl.binary.Float32Serializer, 2)}, {@test_model.AcquisitionOrImage.Acquisition, @test_model.AcquisitionOrImage.Image}));
      w.write(obj.stream_, value);
    end
  end
end