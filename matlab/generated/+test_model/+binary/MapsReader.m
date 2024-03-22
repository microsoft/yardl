% Binary reader for the Maps protocol
classdef MapsReader < yardl.binary.BinaryProtocolReader & test_model.MapsReaderBase
  methods
    function obj = MapsReader(filename)
      obj@test_model.MapsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.MapsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_string_to_int_(obj)
      r = yardl.binary.MapSerializer(yardl.binary.StringSerializer, yardl.binary.Int32Serializer);
      value = r.read(obj.stream_);
    end

    function value = read_int_to_string_(obj)
      r = yardl.binary.MapSerializer(yardl.binary.Int32Serializer, yardl.binary.StringSerializer);
      value = r.read(obj.stream_);
    end

    function value = read_string_to_union_(obj)
      r = yardl.binary.MapSerializer(yardl.binary.StringSerializer, yardl.binary.UnionSerializer('test_model.StringOrInt32', {yardl.binary.StringSerializer, yardl.binary.Int32Serializer}, {@test_model.StringOrInt32.String, @test_model.StringOrInt32.Int32}));
      value = r.read(obj.stream_);
    end

    function value = read_aliased_generic_(obj)
      r = yardl.binary.MapSerializer(yardl.binary.StringSerializer, yardl.binary.Int32Serializer);
      value = r.read(obj.stream_);
    end
  end
end
