% Binary reader for the Enums protocol
classdef EnumsReader < yardl.binary.BinaryProtocolReader & test_model.EnumsReaderBase
  methods
    function obj = EnumsReader(filename)
      obj@test_model.EnumsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.EnumsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_single_(obj)
      r = yardl.binary.EnumSerializer('basic_types.Fruits', @basic_types.Fruits, yardl.binary.Int32Serializer);
      value = r.read(obj.stream_);
    end

    function value = read_vec_(obj)
      r = yardl.binary.VectorSerializer(yardl.binary.EnumSerializer('basic_types.Fruits', @basic_types.Fruits, yardl.binary.Int32Serializer));
      value = r.read(obj.stream_);
    end

    function value = read_size_(obj)
      r = yardl.binary.EnumSerializer('test_model.SizeBasedEnum', @test_model.SizeBasedEnum, yardl.binary.SizeSerializer);
      value = r.read(obj.stream_);
    end
  end
end
