% Binary reader for the NestedRecords protocol
classdef NestedRecordsReader < yardl.binary.BinaryProtocolReader & test_model.NestedRecordsReaderBase
  methods
    function obj = NestedRecordsReader(filename)
      obj@test_model.NestedRecordsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.NestedRecordsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_tuple_with_records_(obj)
      r = test_model.binary.TupleWithRecordsSerializer();
      value = r.read(obj.stream_);
    end
  end
end
