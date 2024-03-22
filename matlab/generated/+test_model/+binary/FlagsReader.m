% Binary reader for the Flags protocol
classdef FlagsReader < yardl.binary.BinaryProtocolReader & test_model.FlagsReaderBase
  methods
    function obj = FlagsReader(filename)
      obj@test_model.FlagsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.FlagsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_days_(obj)
      r = yardl.binary.StreamSerializer(yardl.binary.EnumSerializer('basic_types.DaysOfWeek', @basic_types.DaysOfWeek, yardl.binary.Int32Serializer));
      value = r.read(obj.stream_);
    end

    function value = read_formats_(obj)
      r = yardl.binary.StreamSerializer(yardl.binary.EnumSerializer('basic_types.TextFormat', @basic_types.TextFormat, yardl.binary.Uint64Serializer));
      value = r.read(obj.stream_);
    end
  end
end
