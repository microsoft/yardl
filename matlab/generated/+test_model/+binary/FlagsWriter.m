% Binary writer for the Flags protocol
classdef FlagsWriter < yardl.binary.BinaryProtocolWriter & test_model.FlagsWriterBase
  methods
    function obj = FlagsWriter(filename)
      obj@test_model.FlagsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.FlagsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_days_(obj, value)
      w = yardl.binary.StreamSerializer(yardl.binary.EnumSerializer('basic_types.DaysOfWeek', @basic_types.DaysOfWeek, yardl.binary.Int32Serializer));
      w.write(obj.stream_, value);
    end

    function write_formats_(obj, value)
      w = yardl.binary.StreamSerializer(yardl.binary.EnumSerializer('basic_types.TextFormat', @basic_types.TextFormat, yardl.binary.Uint64Serializer));
      w.write(obj.stream_, value);
    end
  end
end
