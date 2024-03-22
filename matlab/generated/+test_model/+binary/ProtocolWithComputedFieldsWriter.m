% Binary writer for the ProtocolWithComputedFields protocol
classdef ProtocolWithComputedFieldsWriter < yardl.binary.BinaryProtocolWriter & test_model.ProtocolWithComputedFieldsWriterBase
  methods
    function obj = ProtocolWithComputedFieldsWriter(filename)
      obj@test_model.ProtocolWithComputedFieldsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.ProtocolWithComputedFieldsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_record_with_computed_fields_(obj, value)
      w = test_model.binary.RecordWithComputedFieldsSerializer();
      w.write(obj.stream_, value);
    end
  end
end
