% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MapsReader < yardl.binary.BinaryProtocolReader & test_model.MapsReaderBase
  % Binary reader for the Maps protocol
  properties (Access=protected)
    string_to_int_serializer
    int_to_string_serializer
    string_to_union_serializer
    aliased_generic_serializer
  end

  methods
    function obj = MapsReader(filename)
      obj@test_model.MapsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.MapsReaderBase.schema);
      obj.string_to_int_serializer = yardl.binary.MapSerializer(yardl.binary.StringSerializer, yardl.binary.Int32Serializer);
      obj.int_to_string_serializer = yardl.binary.MapSerializer(yardl.binary.Int32Serializer, yardl.binary.StringSerializer);
      obj.string_to_union_serializer = yardl.binary.MapSerializer(yardl.binary.StringSerializer, yardl.binary.UnionSerializer('test_model.StringOrInt32', {yardl.binary.StringSerializer, yardl.binary.Int32Serializer}, {@test_model.StringOrInt32.String, @test_model.StringOrInt32.Int32}));
      obj.aliased_generic_serializer = yardl.binary.MapSerializer(yardl.binary.StringSerializer, yardl.binary.Int32Serializer);
    end
  end

  methods (Access=protected)
    function value = read_string_to_int_(obj)
      value = obj.string_to_int_serializer.read(obj.stream_);
    end

    function value = read_int_to_string_(obj)
      value = obj.int_to_string_serializer.read(obj.stream_);
    end

    function value = read_string_to_union_(obj)
      value = obj.string_to_union_serializer.read(obj.stream_);
    end

    function value = read_aliased_generic_(obj)
      value = obj.aliased_generic_serializer.read(obj.stream_);
    end
  end
end
