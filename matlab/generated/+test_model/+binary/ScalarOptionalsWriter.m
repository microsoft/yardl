% Binary writer for the ScalarOptionals protocol
classdef ScalarOptionalsWriter < yardl.binary.BinaryProtocolWriter & test_model.ScalarOptionalsWriterBase
  methods
    function obj = ScalarOptionalsWriter(filename)
      obj@test_model.ScalarOptionalsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.ScalarOptionalsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_optional_int_(obj, value)
      w = yardl.binary.OptionalSerializer(yardl.binary.Int32Serializer);
      w.write(obj.stream_, value);
    end

    function write_optional_record_(obj, value)
      w = yardl.binary.OptionalSerializer(test_model.binary.SimpleRecordSerializer());
      w.write(obj.stream_, value);
    end

    function write_record_with_optional_fields_(obj, value)
      w = test_model.binary.RecordWithOptionalFieldsSerializer();
      w.write(obj.stream_, value);
    end

    function write_optional_record_with_optional_fields_(obj, value)
      w = yardl.binary.OptionalSerializer(test_model.binary.RecordWithOptionalFieldsSerializer());
      w.write(obj.stream_, value);
    end
  end
end
