% Binary writer for the Enums protocol
classdef EnumsWriter < yardl.binary.BinaryProtocolWriter & test_model.EnumsWriterBase
  methods
    function obj = EnumsWriter(filename)
      obj@test_model.EnumsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.EnumsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_single_(obj, value)
      w = yardl.binary.EnumSerializer('basic_types.Fruits', @basic_types.Fruits, yardl.binary.Int32Serializer);
      w.write(obj.stream_, value);
    end

    function write_vec_(obj, value)
      w = yardl.binary.VectorSerializer(yardl.binary.EnumSerializer('basic_types.Fruits', @basic_types.Fruits, yardl.binary.Int32Serializer));
      w.write(obj.stream_, value);
    end

    function write_size_(obj, value)
      w = yardl.binary.EnumSerializer('test_model.SizeBasedEnum', @test_model.SizeBasedEnum, yardl.binary.SizeSerializer);
      w.write(obj.stream_, value);
    end
  end
end
