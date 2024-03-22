% Binary reader for the ScalarOptionals protocol
classdef ScalarOptionalsReader < yardl.binary.BinaryProtocolReader & test_model.ScalarOptionalsReaderBase
  methods
    function obj = ScalarOptionalsReader(filename)
      obj@test_model.ScalarOptionalsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.ScalarOptionalsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_optional_int_(obj)
      r = yardl.binary.OptionalSerializer(yardl.binary.Int32Serializer);
      value = r.read(obj.stream_);
    end

    function value = read_optional_record_(obj)
      r = yardl.binary.OptionalSerializer(test_model.binary.SimpleRecordSerializer());
      value = r.read(obj.stream_);
    end

    function value = read_record_with_optional_fields_(obj)
      r = test_model.binary.RecordWithOptionalFieldsSerializer();
      value = r.read(obj.stream_);
    end

    function value = read_optional_record_with_optional_fields_(obj)
      r = yardl.binary.OptionalSerializer(test_model.binary.RecordWithOptionalFieldsSerializer());
      value = r.read(obj.stream_);
    end
  end
end