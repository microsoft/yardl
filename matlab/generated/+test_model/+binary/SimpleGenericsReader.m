% Binary reader for the SimpleGenerics protocol
classdef SimpleGenericsReader < yardl.binary.BinaryProtocolReader & test_model.SimpleGenericsReaderBase
  methods
    function obj = SimpleGenericsReader(filename)
      obj@test_model.SimpleGenericsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.SimpleGenericsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_float_image_(obj)
      r = yardl.binary.NDArraySerializer(yardl.binary.Float32Serializer, 2);
      value = r.read(obj.stream_);
    end

    function value = read_int_image_(obj)
      r = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      value = r.read(obj.stream_);
    end

    function value = read_int_image_alternate_syntax_(obj)
      r = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      value = r.read(obj.stream_);
    end

    function value = read_string_image_(obj)
      r = yardl.binary.NDArraySerializer(yardl.binary.StringSerializer, 2);
      value = r.read(obj.stream_);
    end

    function value = read_int_float_tuple_(obj)
      r = tuples.binary.TupleSerializer(yardl.binary.Int32Serializer, yardl.binary.Float32Serializer);
      value = r.read(obj.stream_);
    end

    function value = read_float_float_tuple_(obj)
      r = tuples.binary.TupleSerializer(yardl.binary.Float32Serializer, yardl.binary.Float32Serializer);
      value = r.read(obj.stream_);
    end

    function value = read_int_float_tuple_alternate_syntax_(obj)
      r = tuples.binary.TupleSerializer(yardl.binary.Int32Serializer, yardl.binary.Float32Serializer);
      value = r.read(obj.stream_);
    end

    function value = read_int_string_tuple_(obj)
      r = tuples.binary.TupleSerializer(yardl.binary.Int32Serializer, yardl.binary.StringSerializer);
      value = r.read(obj.stream_);
    end

    function value = read_stream_of_type_variants_(obj)
      r = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('test_model.ImageFloatOrImageDouble', {yardl.binary.NDArraySerializer(yardl.binary.Float32Serializer, 2), yardl.binary.NDArraySerializer(yardl.binary.Float64Serializer, 2)}, {@test_model.ImageFloatOrImageDouble.ImageFloat, @test_model.ImageFloatOrImageDouble.ImageDouble}));
      value = r.read(obj.stream_);
    end
  end
end
