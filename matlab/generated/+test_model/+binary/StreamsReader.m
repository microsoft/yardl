% Binary reader for the Streams protocol
classdef StreamsReader < yardl.binary.BinaryProtocolReader & test_model.StreamsReaderBase
  methods
    function obj = StreamsReader(filename)
      obj@test_model.StreamsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.StreamsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_int_data_(obj)
      r = yardl.binary.StreamSerializer(yardl.binary.Int32Serializer);
      value = r.read(obj.stream_);
    end

    function value = read_optional_int_data_(obj)
      r = yardl.binary.StreamSerializer(yardl.binary.OptionalSerializer(yardl.binary.Int32Serializer));
      value = r.read(obj.stream_);
    end

    function value = read_record_with_optional_vector_data_(obj)
      r = yardl.binary.StreamSerializer(test_model.binary.RecordWithOptionalVectorSerializer());
      value = r.read(obj.stream_);
    end

    function value = read_fixed_vector_(obj)
      r = yardl.binary.StreamSerializer(yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 3));
      value = r.read(obj.stream_);
    end
  end
end