% Binary reader for the AdvancedGenerics protocol
classdef AdvancedGenericsReader < yardl.binary.BinaryProtocolReader & test_model.AdvancedGenericsReaderBase
  methods
    function obj = AdvancedGenericsReader(filename)
      obj@test_model.AdvancedGenericsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.AdvancedGenericsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_float_image_image_(obj)
      r = yardl.binary.NDArraySerializer(yardl.binary.NDArraySerializer(yardl.binary.Float32Serializer, 2), 2);
      value = r.read(obj.stream_);
    end

    function value = read_generic_record_1_(obj)
      r = test_model.binary.GenericRecordSerializer(yardl.binary.Int32Serializer, yardl.binary.StringSerializer);
      value = r.read(obj.stream_);
    end

    function value = read_tuple_of_optionals_(obj)
      r = tuples.binary.TupleSerializer(yardl.binary.OptionalSerializer(yardl.binary.Int32Serializer), yardl.binary.OptionalSerializer(yardl.binary.StringSerializer));
      value = r.read(obj.stream_);
    end

    function value = read_tuple_of_optionals_alternate_syntax_(obj)
      r = tuples.binary.TupleSerializer(yardl.binary.OptionalSerializer(yardl.binary.Int32Serializer), yardl.binary.OptionalSerializer(yardl.binary.StringSerializer));
      value = r.read(obj.stream_);
    end

    function value = read_tuple_of_vectors_(obj)
      r = tuples.binary.TupleSerializer(yardl.binary.VectorSerializer(yardl.binary.Int32Serializer), yardl.binary.VectorSerializer(yardl.binary.Float32Serializer));
      value = r.read(obj.stream_);
    end
  end
end
