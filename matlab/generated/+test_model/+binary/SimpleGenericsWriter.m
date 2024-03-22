% Binary writer for the SimpleGenerics protocol
classdef SimpleGenericsWriter < yardl.binary.BinaryProtocolWriter & test_model.SimpleGenericsWriterBase
  methods
    function obj = SimpleGenericsWriter(filename)
      obj@test_model.SimpleGenericsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.SimpleGenericsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_float_image_(obj, value)
      w = yardl.binary.NDArraySerializer(yardl.binary.Float32Serializer, 2);
      w.write(obj.stream_, value);
    end

    function write_int_image_(obj, value)
      w = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      w.write(obj.stream_, value);
    end

    function write_int_image_alternate_syntax_(obj, value)
      w = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      w.write(obj.stream_, value);
    end

    function write_string_image_(obj, value)
      w = yardl.binary.NDArraySerializer(yardl.binary.StringSerializer, 2);
      w.write(obj.stream_, value);
    end

    function write_int_float_tuple_(obj, value)
      w = tuples.binary.TupleSerializer(yardl.binary.Int32Serializer, yardl.binary.Float32Serializer);
      w.write(obj.stream_, value);
    end

    function write_float_float_tuple_(obj, value)
      w = tuples.binary.TupleSerializer(yardl.binary.Float32Serializer, yardl.binary.Float32Serializer);
      w.write(obj.stream_, value);
    end

    function write_int_float_tuple_alternate_syntax_(obj, value)
      w = tuples.binary.TupleSerializer(yardl.binary.Int32Serializer, yardl.binary.Float32Serializer);
      w.write(obj.stream_, value);
    end

    function write_int_string_tuple_(obj, value)
      w = tuples.binary.TupleSerializer(yardl.binary.Int32Serializer, yardl.binary.StringSerializer);
      w.write(obj.stream_, value);
    end

    function write_stream_of_type_variants_(obj, value)
      w = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('test_model.ImageFloatOrImageDouble', {yardl.binary.NDArraySerializer(yardl.binary.Float32Serializer, 2), yardl.binary.NDArraySerializer(yardl.binary.Float64Serializer, 2)}, {@test_model.ImageFloatOrImageDouble.ImageFloat, @test_model.ImageFloatOrImageDouble.ImageDouble}));
      w.write(obj.stream_, value);
    end
  end
end
