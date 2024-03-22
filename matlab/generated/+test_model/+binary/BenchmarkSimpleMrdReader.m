% Binary reader for the BenchmarkSimpleMrd protocol
classdef BenchmarkSimpleMrdReader < yardl.binary.BinaryProtocolReader & test_model.BenchmarkSimpleMrdReaderBase
  methods
    function obj = BenchmarkSimpleMrdReader(filename)
      obj@test_model.BenchmarkSimpleMrdReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.BenchmarkSimpleMrdReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_data_(obj)
      r = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('test_model.AcquisitionOrImage', {test_model.binary.SimpleAcquisitionSerializer(), yardl.binary.NDArraySerializer(yardl.binary.Float32Serializer, 2)}, {@test_model.AcquisitionOrImage.Acquisition, @test_model.AcquisitionOrImage.Image}));
      value = r.read(obj.stream_);
    end
  end
end
