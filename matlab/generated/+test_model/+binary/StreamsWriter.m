% Binary writer for the Streams protocol
classdef StreamsWriter < yardl.binary.BinaryProtocolWriter & test_model.StreamsWriterBase
  methods
    function obj = StreamsWriter(filename)
      obj@test_model.StreamsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.StreamsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_int_data_(obj, value)
      w = yardl.binary.StreamSerializer(yardl.binary.Int32Serializer);
      w.write(obj.stream_, value);
    end

    function write_optional_int_data_(obj, value)
      w = yardl.binary.StreamSerializer(yardl.binary.OptionalSerializer(yardl.binary.Int32Serializer));
      w.write(obj.stream_, value);
    end

    function write_record_with_optional_vector_data_(obj, value)
      w = yardl.binary.StreamSerializer(test_model.binary.RecordWithOptionalVectorSerializer());
      w.write(obj.stream_, value);
    end

    function write_fixed_vector_(obj, value)
      w = yardl.binary.StreamSerializer(yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 3));
      w.write(obj.stream_, value);
    end
  end
end
