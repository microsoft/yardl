% Binary reader for the ProtocolWithComputedFields protocol
classdef ProtocolWithComputedFieldsReader < yardl.binary.BinaryProtocolReader & test_model.ProtocolWithComputedFieldsReaderBase
  methods
    function obj = ProtocolWithComputedFieldsReader(filename)
      obj@test_model.ProtocolWithComputedFieldsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.ProtocolWithComputedFieldsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_record_with_computed_fields_(obj)
      r = test_model.binary.RecordWithComputedFieldsSerializer();
      value = r.read(obj.stream_);
    end
  end
end
