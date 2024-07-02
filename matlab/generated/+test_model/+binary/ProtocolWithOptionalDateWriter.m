% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef ProtocolWithOptionalDateWriter < yardl.binary.BinaryProtocolWriter & test_model.ProtocolWithOptionalDateWriterBase
  % Binary writer for the ProtocolWithOptionalDate protocol
  properties (Access=protected)
    record_serializer
  end

  methods
    function self = ProtocolWithOptionalDateWriter(filename)
      self@test_model.ProtocolWithOptionalDateWriterBase();
      self@yardl.binary.BinaryProtocolWriter(filename, test_model.ProtocolWithOptionalDateWriterBase.schema);
      self.record_serializer = yardl.binary.OptionalSerializer(test_model.binary.RecordWithOptionalDateSerializer());
    end
  end

  methods (Access=protected)
    function write_record_(self, value)
      self.record_serializer.write(self.stream_, value);
    end
  end
end
