% Binary writer for the AdvancedGenerics protocol
classdef AdvancedGenericsWriter < yardl.binary.BinaryProtocolWriter & test_model.AdvancedGenericsWriterBase
  methods
    function obj = AdvancedGenericsWriter(filename)
      obj@test_model.AdvancedGenericsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.AdvancedGenericsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_float_image_image_(obj, value)
      w = yardl.binary.NDArraySerializer(yardl.binary.NDArraySerializer(yardl.binary.Float32Serializer, 2), 2);
      w.write(obj.stream_, value);
    end

    function write_generic_record_1_(obj, value)
      w = test_model.binary.GenericRecordSerializer(yardl.binary.Int32Serializer, yardl.binary.StringSerializer);
      w.write(obj.stream_, value);
    end

    function write_tuple_of_optionals_(obj, value)
      w = tuples.binary.TupleSerializer(yardl.binary.OptionalSerializer(yardl.binary.Int32Serializer), yardl.binary.OptionalSerializer(yardl.binary.StringSerializer));
      w.write(obj.stream_, value);
    end

    function write_tuple_of_optionals_alternate_syntax_(obj, value)
      w = tuples.binary.TupleSerializer(yardl.binary.OptionalSerializer(yardl.binary.Int32Serializer), yardl.binary.OptionalSerializer(yardl.binary.StringSerializer));
      w.write(obj.stream_, value);
    end

    function write_tuple_of_vectors_(obj, value)
      w = tuples.binary.TupleSerializer(yardl.binary.VectorSerializer(yardl.binary.Int32Serializer), yardl.binary.VectorSerializer(yardl.binary.Float32Serializer));
      w.write(obj.stream_, value);
    end
  end
end
