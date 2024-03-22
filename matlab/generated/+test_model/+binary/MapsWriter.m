% Binary writer for the Maps protocol
classdef MapsWriter < yardl.binary.BinaryProtocolWriter & test_model.MapsWriterBase
  methods
    function obj = MapsWriter(filename)
      obj@test_model.MapsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.MapsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_string_to_int_(obj, value)
      w = yardl.binary.MapSerializer(yardl.binary.StringSerializer, yardl.binary.Int32Serializer);
      w.write(obj.stream_, value);
    end

    function write_int_to_string_(obj, value)
      w = yardl.binary.MapSerializer(yardl.binary.Int32Serializer, yardl.binary.StringSerializer);
      w.write(obj.stream_, value);
    end

    function write_string_to_union_(obj, value)
      w = yardl.binary.MapSerializer(yardl.binary.StringSerializer, yardl.binary.UnionSerializer('test_model.StringOrInt32', {yardl.binary.StringSerializer, yardl.binary.Int32Serializer}, {@test_model.StringOrInt32.String, @test_model.StringOrInt32.Int32}));
      w.write(obj.stream_, value);
    end

    function write_aliased_generic_(obj, value)
      w = yardl.binary.MapSerializer(yardl.binary.StringSerializer, yardl.binary.Int32Serializer);
      w.write(obj.stream_, value);
    end
  end
end
