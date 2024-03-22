% Binary writer for the NestedRecords protocol
classdef NestedRecordsWriter < yardl.binary.BinaryProtocolWriter & test_model.NestedRecordsWriterBase
  methods
    function obj = NestedRecordsWriter(filename)
      obj@test_model.NestedRecordsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.NestedRecordsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_tuple_with_records_(obj, value)
      w = test_model.binary.TupleWithRecordsSerializer();
      w.write(obj.stream_, value);
    end
  end
end
